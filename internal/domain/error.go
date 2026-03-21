package domain



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

type ErrorDetail struct {
    Code    string `json:"code"`
    Message string `json:"message"`
}

type ErrorResponse struct {
    Error ErrorDetail `json:"error"`
}

type InternalErrorResponse struct {
    Error ErrorDetail `json:"error"`
}

func NewError(code, message string) ErrorResponse {
    return ErrorResponse{
        Error: ErrorDetail{
            Code:    code,
            Message: message,
        },
    }
}