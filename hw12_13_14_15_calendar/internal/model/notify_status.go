package model

import "errors"

var ErrorUnknownNotifyStatus = errors.New("unknown notify status")

type NotifyStatus int

const (
	NotifyStatusNone NotifyStatus = iota
	NotifyStatusBlocked
	NotifyStatusNotified
	NotifyStatusError
)

func (ns NotifyStatus) Valid() bool {
	return ns < NotifyStatusError
}

func (ns NotifyStatus) String() string {
	switch ns {
	case NotifyStatusBlocked:
		return "blocked"
	case NotifyStatusNotified:
		return "notified"
	case NotifyStatusNone:
		return "none"
	}
	return ""
}

func ParseNotifyStatus(status string) (NotifyStatus, error) {
	switch status {
	case "none":
		return NotifyStatusNone, nil
	case "blocked":
		return NotifyStatusBlocked, nil
	case "notified":
		return NotifyStatusNotified, nil
	}
	return NotifyStatusError, ErrorUnknownNotifyStatus
}
