package model

import "errors"

var (
	ErrUserEmptyID        = errors.New("не задан ID пользователя")
	ErrUserEmptyName      = errors.New("не введено имя пользователя")
	ErrUserWrongEmail     = errors.New("неверный E-mail")
	ErrUserDuplicateEmail = errors.New("пользователь с таким E-mail уже существует")
	ErrUserNotFound       = errors.New("указанный пользователь не найден")
)
