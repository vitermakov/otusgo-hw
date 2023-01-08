package model

import (
	"github.com/google/uuid"
	"time"
)

// Notification модель оповещения о приближающемся событии.
type Notification struct {
	EventID       uuid.UUID     `json:"eventID"`       // ID события.
	EventTitle    string        `json:"eventTitle"`    // Заголовок события.
	EventDate     time.Time     `json:"eventDate"`     // Дата события.
	EventDuration time.Duration `json:"eventDuration"` // Продолжительность события.
	NotifyUser    NotifyUser    `json:"notifyUser"`    // Пользователь, которому отправлять.
}

type NotifyUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
