package model

import "errors"

var (
	ErrEventEmptyTitle   = errors.New("заголовок пуст")
	ErrEventZeroDate     = errors.New("неверная дата начала")
	ErrEventZeroDuration = errors.New("неверная длительность")
	ErrEventOwnerId      = errors.New("не задан идентификатор владельца")
	ErrEventOwnerExists  = errors.New("указанный владелец не найден")
	ErrEventDateBusy     = errors.New("указанная дата занята")
	ErrEventNotFound     = errors.New("указанное событие не найдено")
)
