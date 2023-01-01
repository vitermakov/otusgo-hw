package model

import "errors"

var (
	ErrEventEmptyTitle    = errors.New("заголовок пуст")
	ErrEventZeroDate      = errors.New("неверная дата начала")
	ErrEventWrongDuration = errors.New("неверная длительность")
)

/*
Если в результате сервисной фукции возвращается ошибка бизнес-логики errx.Logic,
то в поле `code` будет вписан целочисленных код внутренней классификации ошибок.
Это упростит тестирование REST.API, так как errx.Logic переходит в ответ rqres.BadRequest, и
проверку можно осущетсвлять по коду
*/

const (
	ErrCalendarAccessCode    = 1001
	ErrCalendarDateRangeCode = 1002
	ErrEventOwnerIDCode      = 1003
	ErrEventOwnerExistsCode  = 1004
	ErrEventDateBusyCode     = 1005
)

var (
	ErrCalendarAccess       = errors.New("нет доступа к календарю")
	ErrCalendarDateRange    = errors.New("неверный интервал дат для получения событий")
	ErrEventOwnerID         = errors.New("не задан идентификатор владельца")
	ErrEventOwnerExists     = errors.New("указанный владелец не найден")
	ErrEventDateBusy        = errors.New("указанная дата занята")
	ErrEventNotFound        = errors.New("указанное событие не найдено")
	ErrEventWrongNotifyTerm = errors.New("неверный интервал оповещения")
	ErrEventNotFoundID      = errors.New("не найдено событие")
)
