package service

import (
	"context"
	domain "test-backend-1-curboturbo/internal/domain"
	port "test-backend-1-curboturbo/internal/port/outbound"
)


type RoomService interface{
	DisplayRooms(ctx context.Context) ([]domain.Room, error)
	CreateRoom(ctx context.Context, req domain.RoomCreateRequest) (domain.Room, error)
}

type roomService struct{
	storage port.StorageProvider
}

func NewRoomService(storage port.StorageProvider) RoomService{
	return &roomService{storage: storage}
}

func (s *roomService) DisplayRooms(ctx context.Context) ([]domain.Room, error) {
	rooms, err := s.storage.ShowRooms(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]domain.Room, 0, len(rooms))
	for _, r := range rooms {
		result = append(result, domain.Room{
			ID:          r.ID,
			Name:        r.Name,
			Description: r.Description,
			Capacity:    r.Capacity,
			CreatedAt:   r.CreatedAt,
		})
	}
	return result, nil
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
	return domain.Room{
		ID: createdRoom.ID,
		Capacity: createdRoom.Capacity,
		CreatedAt: createdRoom.CreatedAt,
		Description: createdRoom.Description,
		Name: createdRoom.Name,
	}, nil
}