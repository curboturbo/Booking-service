package storage

import (
	"context"
	"errors"
	"fmt"
	"strings"
	domain "test-backend-1-curboturbo/internal/domain"
	models "test-backend-1-curboturbo/internal/model"
	port "test-backend-1-curboturbo/internal/port/outbound"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type storageConnector struct {
    db *gorm.DB
}


func NewStorage(db *gorm.DB) port.StorageProvider {
	err := db.AutoMigrate(
    &models.User{},
    &models.Room{},
    &models.Schedule{},
    &models.Slot{},
    &models.Booking{},
    )
    if err != nil{
        panic("cannot transfer migrations")
    }
    return &storageConnector{
        db: db,
    }
}

func ToPQInt32Array(in []int) pq.Int32Array {
	out := make(pq.Int32Array, len(in))
	for i, v := range in {
		out[i] = int32(v)
	}
	return out
}


func (s *storageConnector) Create(ctx context.Context, email string, password string, role string) (models.User, error) {
    user := models.User{
        Email:    email,
        Password: password,
        Role:     role,
    }
    result := s.db.WithContext(ctx).Create(&user)
    if result.Error != nil {
        if strings.Contains(result.Error.Error(), "duplicate key") || 
           strings.Contains(result.Error.Error(), "UNIQUE constraint") {
            return models.User{}, domain.ErrEmailAlreadyTaken
        }
        return models.User{}, result.Error
    }
    return user, nil
}

func (s* storageConnector) CreateAdmin(ctx context.Context, email string, password string) (error){
    user := models.User{Email: email}
    result := s.db.WithContext(ctx).
        Where(models.User{Email: email}).
        FirstOrCreate(&user, models.User{Password: password, Role:"admin"})
    if result.Error != nil{return result.Error}
    return nil
}


func (s *storageConnector) GetUser(ctx context.Context, email string) (uuid.UUID, string, string, error) {
    var user models.User
    result := s.db.WithContext(ctx).
        Select("id", "password", "role"). 
        Where("email = ?", email).
        First(&user)
    if result.Error != nil {
        if errors.Is(result.Error, gorm.ErrRecordNotFound) {
            return uuid.Nil, "", "", domain.ErrUserNotFound
        }
        return uuid.Nil, "", "", result.Error
    }
    return user.ID, user.Password, user.Role, nil
}


func (s *storageConnector) ShowRooms(ctx context.Context) ([]domain.Room, error){
    var rooms []models.Room
    result := s.db.WithContext(ctx).Find(&rooms)
    if result.Error != nil{
        return []domain.Room{}, result.Error
    }
    res := make([]domain.Room, 0, len(rooms))
	for _, r := range rooms {
		res = append(res, domain.Room{
			ID:          r.ID,
			Name:        r.Name,
			Description: r.Description,
			Capacity:    r.Capacity,
			CreatedAt:   r.CreatedAt,
		})
	}
    return res, nil
}


func (s *storageConnector) CreateRoom(ctx context.Context, room domain.Room) (domain.Room, error) {
	model := models.Room{
		Name:        room.Name,
		Description: room.Description,
		Capacity:    room.Capacity,
	}
	result := s.db.WithContext(ctx).Create(&model)
	if result.Error != nil {
		return domain.Room{}, result.Error
	}

	return domain.Room{
		Capacity: model.Capacity,
		CreatedAt: model.CreatedAt,
		Description: model.Description,
		Name: model.Name,
	}, nil
}


func (s *storageConnector) CreateSchedule(ctx context.Context, sched domain.Schedule) (domain.Schedule, error){
    schedule := models.Schedule{
        RoomID: sched.RoomID,
        DaysOfWeek: ToPQInt32Array(sched.DaysOfWeek),
        StartTime: sched.StartTime,
        EndTime: sched.EndTime,
    }
    fmt.Print(schedule.RoomID)
    var pqErr *pq.Error
    var room models.Room
    err := s.db.First(&room, "id = ?", schedule.RoomID).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return domain.Schedule{}, domain.RoomNotFound
        }
        return domain.Schedule{}, domain.InternalError
    }
    res := s.db.WithContext(ctx).Create(&schedule)

    if res.Error != nil {
        if errors.As(res.Error, &pqErr) && pqErr.Code == "23505" {
            return domain.Schedule{}, domain.ErrSchedultAlreayExist
        }
        return domain.Schedule{}, domain.InternalError
    }

    var slots []models.Slot
    const delta = 30*time.Minute
    timeLayout := "15:04"
    slotBegin , _ := time.Parse(timeLayout, sched.StartTime)
    slotEnd, _ := time.Parse(timeLayout, sched.EndTime)
    now := time.Now().UTC()
    cleanDay := time.Date(
        now.Year(),
        now.Month(),
        now.Day(),
        0,0,0,0,
        time.UTC,
    )
    for i:=0;i<=14;i++{
        curDate := cleanDay.AddDate(0,0,i)
        day := int(curDate.Weekday())
        if day == 0 {
            day = 7
        }
        exist := false
        for j:=range sched.DaysOfWeek{
            if sched.DaysOfWeek[j] == day{
                exist = true
                break
            }
        }
        if !exist{continue}

        slotStart := time.Date(
            curDate.Year(), curDate.Month(), curDate.Day(),
            slotBegin.Hour(), slotBegin.Minute(), 0, 0, time.UTC,
        )
        
        dayEndLimit := time.Date(
            curDate.Year(), curDate.Month(), curDate.Day(),
            slotEnd.Hour(), slotEnd.Minute(), 0, 0, time.UTC,
        )
        for{
            nextWindow := slotStart.Add(delta)
            if nextWindow.After(dayEndLimit){
                break
            }
            slots = append(
                slots,
                models.Slot{
                    RoomID: schedule.RoomID,
                    StartTime: slotStart,
                    EndTime: nextWindow,
                },
            )
            slotStart = nextWindow
        }
    }
    result := s.db.WithContext(ctx).CreateInBatches(slots, len(slots)-1)
    if result.Error != nil{
        return domain.Schedule{}, domain.InternalError
    }
    return domain.Schedule{
		RoomID: schedule.RoomID,
		StartTime: schedule.StartTime,
		EndTime: schedule.EndTime,
	}, nil
}

func (s *storageConnector) TakeSlots(ctx context.Context, filterSlot domain.Slot) ([]domain.Slot, error){
    var availableslots []models.Slot
    date := time.Date(filterSlot.StartTime.Year(), filterSlot.StartTime.Month(), filterSlot.StartTime.Day(), 0, 0, 0, 0, time.UTC)
    err := s.db.WithContext(ctx).
        Select("id, room_id, start_time, end_time").
        Where("room_id = ? AND start_time < ?", filterSlot.RoomId, date).
        Where("NOT EXISTS (SELECT 1 FROM bookings WHERE bookings.slot_id = slots.id AND bookings.status = 'active')").
        Find(&availableslots).Error
    if err != nil{
        return []domain.Slot{}, domain.InternalError
    }
    if len(availableslots) == 0{return []domain.Slot{}, nil}
    slots := make([]domain.Slot,len(availableslots))
    for i := range availableslots{
        slots[i].RoomId = availableslots[i].RoomID
        slots[i].StartTime = availableslots[i].StartTime
        slots[i].EndTime = availableslots[i].EndTime
    }
    return slots, nil
}


func (s *storageConnector) CreateBooking(ctx context.Context, booking domain.Booking) (domain.Booking, error){
    query := `INSERT INTO bookings (slot_id, user_id, status,link,created_at)
    VALUES ($1, $2, 'active', $3, NOW())
    ON CONFLICT (slot_id)
    DO UPDATE SET
        user_id = EXCLUDED.user_id,
        status = 'active',
        link = EXCLUDED.link,
        created_at = NOW()
    WHERE bookings.status = 'cancel'
    RETURNING id, slot_id, user_id, status, link, created_at;
    `
    var potential_booking models.Booking
    res := s.db.WithContext(ctx).Raw(query, booking.SlotID, booking.UserID, booking.Link).Scan(&potential_booking)
    if res.Error != nil{
        return domain.Booking{},domain.InternalError
    }
    if res.RowsAffected == 0{
        return domain.Booking{}, domain.ErrSlotAlreadyTaken
    }
    return domain.Booking{
        ID:potential_booking.ID,
        SlotID:potential_booking.SlotID,
        UserID: potential_booking.UserID,
        Status: potential_booking.Status,
        Link:potential_booking.Link,
    }, nil
}
