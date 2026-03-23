package service

import (
	"context"
	domain "test-backend-1-curboturbo/internal/domain"
	port "test-backend-1-curboturbo/internal/port/outbound"
	"time"
)


type RoomService interface{
	DisplayRooms(ctx context.Context) ([]domain.Room, error)
	CreateRoom(ctx context.Context, req domain.RoomCreateRequest) (domain.Room, error)
	CreateSchedule(ctx context.Context, req domain.ScheduleCreateRequest) (domain.Schedule, error)
	TakeSlots(ctx context.Context, req domain.AvailableSlotRequest) ([]domain.Slot, error)
	ReserveSlot(ctx context.Context, req domain.CreateBookingRequest) (domain.Booking, error)
}

type roomService struct{
	storage port.StorageProvider
	conference port.LinkConferenceProvider
}

func NewRoomService(storage port.StorageProvider) RoomService{
	return &roomService{storage: storage}
}

func (s *roomService) DisplayRooms(ctx context.Context) ([]domain.Room, error) {
	rooms, err := s.storage.ShowRooms(ctx)
	if err != nil {return nil, err}
	return rooms, nil
}


func (s *roomService) CreateRoom(ctx context.Context, req domain.RoomCreateRequest) (domain.Room, error) {
    room := domain.Room{
        Name: req.Name,
        Description: req.Description,
        Capacity: req.Capacity,
    }
    createdRoom,err := s.storage.CreateRoom(ctx, room)
	if err !=nil{
		return domain.Room{}, err
	}
	return createdRoom, err
}


func (s *roomService) CreateSchedule(ctx context.Context, req domain.ScheduleCreateRequest) (domain.Schedule, error){
	schedule := domain.Schedule{
		RoomID: req.RoomID,
		DaysOfWeek: req.DaysOfWeek,
		StartTime: req.StartTime,
		EndTime: req.EndTime,
	}

	createdSchedule,err := s.storage.CreateSchedule(ctx,schedule)
	if err != nil{
		return domain.Schedule{}, err
	}
	
	return domain.Schedule{
		RoomID: createdSchedule.RoomID,
		StartTime: createdSchedule.StartTime,
		EndTime: createdSchedule.EndTime,
	}, nil
}

func (s *roomService) TakeSlots(ctx context.Context, req domain.AvailableSlotRequest) ([]domain.Slot, error){
	slot := domain.Slot{
		RoomId: req.RoomID,
		StartTime: req.Date,
	}
	slots, err := s.storage.TakeSlots(ctx, slot)
	if err != nil{
		return []domain.Slot{}, err
	}
	return slots, nil
}

func (s *roomService) ReserveSlot(ctx context.Context, req domain.CreateBookingRequest) (domain.Booking, error) {
    var conferenceURL string
    if req.CreateConferenceLink {
        apiCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
        defer cancel()
        conferenceURL, _ = s.conference.RequestLink(apiCtx)
    }
    booking := domain.Booking{
        SlotID: req.SlotID,
        UserID: req.UserID,
        Link:   conferenceURL,
    }
    createdBooking, err := s.storage.CreateBooking(ctx, booking)
    if err != nil {
        return domain.Booking{}, err
    }
    return createdBooking, nil
}