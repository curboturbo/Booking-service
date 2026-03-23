package domain
import "errors"



const (
    ErrCodeInvalidRequest      = "INVALID_REQUEST"
    ErrCodeUnauthorized        = "UNAUTHORIZED"
    ErrCodeNotFound            = "NOT_FOUND"
    ErrCodeRoomNotFound        = "ROOM_NOT_FOUND"
    ErrCodeSlotNotFound        = "SLOT_NOT_FOUND"
    ErrCodeSlotAlreadyBooked   = "SLOT_ALREADY_BOOKED"
    ErrCodeBookingNotFound     = "BOOKING_NOT_FOUND"
    ErrCodeForbidden           = "FORBIDDEN"
    ErrCodeScheduleExists      = "SCHEDULE_EXISTS"
    ErrCodeInternalError       = "INTERNAL_ERROR"
)

var ErrEmailAlreadyTaken = errors.New("Пользователь с таким email уже существует")
var ErrSchedultAlreayExist = errors.New("Расписание для переговорки уже создано, изменение не допускается")
var ErrSchedlultNotFound = errors.New("Переговорка не найдена")
var RoomNotFound = errors.New("Комната не найдена")
var InternalError = errors.New("Внутрення ошибка сервера")
var ErrUserNotFound = errors.New("Неверные учетные данные")
var ErrInvalidScheduleData = errors.New("Неверный запрос (в т.ч. недопустимые значения daysOfWeek)")
var ErrInvalidSlotsData = errors.New("Неверный запрос (отсутствует или некорректен параметр date)")
var ErrInvalidCreateBookingData = errors.New("Неверный запрос")

var ErrSlotAlreadyTaken = errors.New("Слот уже занят")
var ErrSlotDosntExist = errors.New("Слот не найден")



type ErrorDetail struct {
    Code    string `json:"code"`
    Message string `json:"message"`
}

type ErrorResponse struct {
    Code    string `json:"code"`
    Message string `json:"message"`
}

type InternalErrorResponse struct {
    Error ErrorDetail `json:"error"`
}

func NewError(code, message string) ErrorResponse {
    return ErrorResponse{
        Code:    code,
        Message: message,
    }
}