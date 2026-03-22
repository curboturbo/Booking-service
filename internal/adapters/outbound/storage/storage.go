package storage

import (
	"context"
	"errors"
	models "test-backend-1-curboturbo/internal/model"
	port "test-backend-1-curboturbo/internal/port/outbound"
    domain "test-backend-1-curboturbo/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
    "strings"
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
            return uuid.Nil, "", "", errors.New("user not found")
        }
        return uuid.Nil, "", "", result.Error
    }
    return user.ID, user.Password, user.Role, nil
}


func (s *storageConnector) ShowRooms(ctx context.Context) ([]models.Room, error){
    var rooms []models.Room
    result := s.db.WithContext(ctx).Find(&rooms)
    if result.Error != nil{
        return []models.Room{}, result.Error
    }
    return rooms, nil
}


func (s *storageConnector) CreateRoom(ctx context.Context, room domain.Room) (models.Room, error) {
	model := models.Room{
		ID:          room.ID,
		Name:        room.Name,
		Description: room.Description,
		Capacity:    room.Capacity,
	}
	err := s.db.WithContext(ctx).Create(&model).Error
	if err != nil {
		return models.Room{}, err
	}
	return model, nil
}