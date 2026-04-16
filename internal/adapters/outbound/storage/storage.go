package storage

import (
	"context"
	"errors"
	"os"
	"strconv"
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
    query := `
        CREATE UNIQUE INDEX IF NOT EXISTS idx_active_bookings_partial 
        ON bookings (slot_id) 
        WHERE status = 'active';`
	if err := db.Exec(query).Error;err!=nil{
		panic("enable create index")
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
        ID: model.ID,
		Capacity: model.Capacity,
		CreatedAt: model.CreatedAt,
		Description: model.Description,
		Name: model.Name,
	}, nil
}


func (s *storageConnector) CreateSchedule(ctx context.Context, sched domain.Schedule) (domain.Schedule, error) {
	var createdSchedule domain.Schedule
	err := s.db.Transaction(func(tx *gorm.DB) error {
		schedule := models.Schedule{
			RoomID:     sched.RoomID,
			DaysOfWeek: ToPQInt32Array(sched.DaysOfWeek),
			StartTime:  sched.StartTime,
			EndTime:    sched.EndTime,
		}

		var room models.Room
		err := tx.WithContext(ctx).First(&room, "id = ?", schedule.RoomID).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domain.RoomNotFound
			}
			return domain.InternalError
		}

		res := tx.WithContext(ctx).Create(&schedule)
		if res.Error != nil {
			var pqErr *pq.Error
			if errors.As(res.Error, &pqErr) && pqErr.Code == "23505" {
				return domain.ErrSchedultAlreayExist
			}
			return domain.InternalError
		}

		var slots []models.Slot
        d, _:= strconv.Atoi(os.Getenv("DELTA"))
		var delta = time.Duration(d) * time.Minute
		timeLayout := "15:04"
		slotBegin, _ := time.Parse(timeLayout, sched.StartTime)
		slotEnd, _ := time.Parse(timeLayout, sched.EndTime)
		now := time.Now().UTC()
		cleanDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

		for i := 0; i <= 14; i++ {
			curDate := cleanDay.AddDate(0, 0, i)
			day := int(curDate.Weekday())
			if day == 0 {
				day = 7
			}
			exist := false
			for j := range sched.DaysOfWeek {
				if sched.DaysOfWeek[j] == day {
					exist = true
					break
				}
			}
			if !exist {
				continue
			}

			slotStart := time.Date(curDate.Year(), curDate.Month(), curDate.Day(), slotBegin.Hour(), slotBegin.Minute(), 0, 0, time.UTC)
			dayEndLimit := time.Date(curDate.Year(), curDate.Month(), curDate.Day(), slotEnd.Hour(), slotEnd.Minute(), 0, 0, time.UTC)

			for {
				nextWindow := slotStart.Add(delta)
				if nextWindow.After(dayEndLimit) {
					break
				}
				slots = append(slots, models.Slot{
					RoomID:    schedule.RoomID,
					StartTime: slotStart,
					EndTime:   nextWindow,
				})
				slotStart = nextWindow
			}
		}
		if len(slots) > 0 {
			result := tx.WithContext(ctx).CreateInBatches(slots, len(slots))
			if result.Error != nil {
				return domain.InternalError
			}
		}
		createdSchedule = domain.Schedule{
			ID:        schedule.ID,
			RoomID:    schedule.RoomID,
			StartTime: schedule.StartTime,
			EndTime:   schedule.EndTime,
		}
		return nil
	})

	if err != nil {
		return domain.Schedule{}, err
	}
	return createdSchedule, nil
}

func (s *storageConnector) TakeSlots(ctx context.Context, filterSlot domain.Slot) ([]domain.Slot, error){
    var availableslots []models.Slot
    date := time.Date(filterSlot.StartTime.Year(), filterSlot.StartTime.Month(), filterSlot.StartTime.Day(), 0, 0, 0, 0, time.UTC)
    err := s.db.WithContext(ctx).
        Select("id, room_id, start_time, end_time").
        Where("room_id = ? AND start_time > ?", filterSlot.RoomId, date).
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



func isUniqueViolation(err error) bool {
    if pgErr, ok := err.(*pq.Error); ok {return pgErr.Code == "23505"}
    return false
}


func (s *storageConnector) CreateBooking(ctx context.Context, booking domain.Booking) (domain.Booking, error){
    query :=`
    INSERT INTO bookings (slot_id, user_id, status, link, created_at)
    VALUES ($1, $2, 'active', $3, NOW())
    RETURNING id, slot_id, user_id, status, link, created_at;`

    var potential_booking models.Booking
    res := s.db.WithContext(ctx).Raw(query, booking.SlotID, booking.UserID, booking.Link).Scan(&potential_booking)
    if res.Error != nil {
        if isUniqueViolation(res.Error) {
            return domain.Booking{}, domain.ErrSlotAlreadyTaken
        }
        return domain.Booking{}, domain.InternalError
    }
    return domain.Booking{
        ID:potential_booking.ID,
        SlotID:potential_booking.SlotID,
        UserID: potential_booking.UserID,
        Status: potential_booking.Status,
        Link:potential_booking.Link,
    }, nil
}

func (s *storageConnector) TakeUserBooking(ctx context.Context, userID uuid.UUID) ([]domain.Booking, error) {
    var user_bookings []models.Booking
    res := s.db.WithContext(ctx).
        Joins("JOIN slots ON bookings.slot_id = slots.id").
        Where("bookings.user_id = ?", userID).
        Where("bookings.status = 'active'").
        Where("slots.start_time > ?", time.Now().UTC()).
        Find(&user_bookings)

    if res.Error != nil {
        return nil, domain.InternalError
    }

    taken_bookings := make([]domain.Booking, len(user_bookings))
    for i := range user_bookings {
        taken_bookings[i].ID = user_bookings[i].ID
        taken_bookings[i].UserID = user_bookings[i].UserID
        taken_bookings[i].SlotID = user_bookings[i].SlotID
        taken_bookings[i].Link = user_bookings[i].Link
        taken_bookings[i].Status = user_bookings[i].Status
    }
    return taken_bookings, nil
}


func (s *storageConnector) CancelUserBooking(ctx context.Context, booking domain.Booking) (domain.Booking, error){
    var cancelled_booking models.Booking
    query := `
    UPDATE bookings
    SET status = 'cancelled'
    WHERE id = ? AND user_id = ?
    RETURNING id, slot_id, user_id, status, link, created_at;
    `
    res := s.db.WithContext(ctx).Raw(query, booking.ID, booking.UserID).Scan(&cancelled_booking)
    if res.Error != nil{
        return domain.Booking{}, domain.InternalError
    }
    if res.RowsAffected == 0{
        var exist int
        s.db.Raw("SELECT 1 FROM bookings WHERE id = ?", booking.ID).Scan(&exist)
        if exist == 0{
            return domain.Booking{}, domain.ErrBookingNotFound
        }
        return domain.Booking{}, domain.ErrTryChandeForeignBooking
    }
    return domain.Booking{
        ID:     cancelled_booking.ID,
        SlotID: cancelled_booking.SlotID,
        UserID: cancelled_booking.UserID,
        Status: cancelled_booking.Status,
        Link:   cancelled_booking.Link,
    }, nil
}


func (s *storageConnector) GetAllBooking(ctx context.Context, req domain.PaginationParams) ([]domain.Booking, error){
    offset := (req.Page - 1) * req.PageSize
    var bookings []models.Booking
    res := s.db.Limit(req.PageSize).Offset(offset).Order("id asc").Find(&bookings)
    if res.Error != nil{
        return []domain.Booking{}, domain.InternalError
    }
    if len(bookings) == 0{return []domain.Booking{},nil}
    taken_bookings := make([]domain.Booking, len(bookings))
    for i := range bookings{
        taken_bookings[i].ID = bookings[i].ID
        taken_bookings[i].SlotID = bookings[i].SlotID
        taken_bookings[i].UserID = bookings[i].UserID
        taken_bookings[i].Status = bookings[i].Status
        taken_bookings[i].Link = bookings[i].Link
    }
    return taken_bookings, nil
}