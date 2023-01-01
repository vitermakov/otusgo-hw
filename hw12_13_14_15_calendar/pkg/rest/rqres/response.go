package rqres

import (
	"errors"
	"net/http"

	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/utils/errx"
)

// Response интерфейс для связи с веб-сервером.
type Response interface {
	GetHTTPCode() int         // HTTP-код ответа
	GetHTTPResp() interface{} // содержимое ответа
	Message() string
	Success() bool
}

// Base общий вид ответа сервера.
type Base struct {
	success bool
	code    int
	message string
	data    interface{}
}

func (res Base) GetHTTPCode() int {
	if res.success {
		return http.StatusOK
	}
	return http.StatusBadRequest
}

func (res Base) GetHTTPResp() interface{} {
	resp := &struct {
		Status  string      `json:"status"`
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data,omitempty"`
	}{
		Code:    res.code,
		Message: res.message,
		Data:    res.data,
	}
	if res.success {
		resp.Status = "success"
	} else {
		resp.Status = "error"
	}
	return resp
}

func (res Base) Success() bool {
	return res.success
}

func (res Base) Message() string {
	return res.message
}

func (res *Base) SetData(data interface{}) *Base {
	res.data = data
	return res
}

// OKResp Удачный ответ на запросы (команды) типа POST/PUT/PATCH/DELETE HTTPCode = 200.
type OKResp struct {
	Base
}

func (res OKResp) GetHTTPCode() int {
	return http.StatusOK
}

func OK(message string, data interface{}) *OKResp {
	if message == "" {
		message = "Успешно"
	}
	return &OKResp{Base{true, http.StatusOK, message, data}}
}

// DataResp Удачный ответ на запрос GET на получение объекта HTTPCode = 200.
type DataResp struct {
	data interface{}
}

func (res DataResp) GetHTTPCode() int {
	return http.StatusOK
}

func (res DataResp) Success() bool {
	return true
}

func (res DataResp) Message() string {
	return ""
}

func (res DataResp) GetHTTPResp() interface{} {
	return res.data
}

func Data(data interface{}) *DataResp {
	return &DataResp{data}
}

// ListResp Удачный ответ на запрос GET на получение списка объектов без навигации HTTPCode = 200.
type ListResp struct {
	items []interface{}
}

func (res ListResp) GetHTTPCode() int {
	return http.StatusOK
}

func (res ListResp) Success() bool {
	return true
}

func (res ListResp) Message() string {
	return ""
}

func (res ListResp) GetHTTPResp() interface{} {
	return res.items
}

func List(list []interface{}) *ListResp {
	return &ListResp{list}
}

// BadResp Ошибка из-за нарушения правил бизнес-логики HTTPCode = 400.
// Ошибки, связанные с действиями пользователей, которые не могут быть выполнены при текущих правилах.
type BadResp struct {
	Base
}

func BadRequest(message string, code int) *BadResp {
	return &BadResp{Base{false, code, message, nil}}
}

// UnAuthResp Попытка совершить операцию без авторизации HTTPCode = 401.
type UnAuthResp struct {
	Base
}

func (res UnAuthResp) GetHTTPCode() int {
	return http.StatusUnauthorized
}

func UnAuth() *UnAuthResp {
	message := "Вы не авторизованы"
	return &UnAuthResp{
		Base{false, http.StatusUnauthorized, message, nil},
	}
}

// NotFoundResp Запрошенный ресурс не найден HTTPCode = 404.
type NotFoundResp struct {
	Base
}

func (res NotFoundResp) GetHTTPCode() int {
	return http.StatusNotFound
}

func NotFound(message string, params interface{}) *NotFoundResp {
	if message == "" {
		message = "Объект не найден"
	}
	return &NotFoundResp{
		Base{false, http.StatusNotFound, message, params},
	}
}

// InvalidResp Ошибка валидации входных данных HTTPCode = 422.
type InvalidResp struct {
	Base
}

func (res InvalidResp) GetHTTPCode() int {
	return http.StatusUnprocessableEntity
}

func (res InvalidResp) GetHTTPResp() interface{} {
	resp := &struct {
		Status  string            `json:"status"`
		Code    int               `json:"code"`
		Message string            `json:"message"`
		Errors  map[string]string `json:"errors"`
	}{
		Status:  "error",
		Code:    res.code,
		Message: res.message,
	}
	if res.data != nil {
		resp.Errors = make(map[string]string)
		errs, ok := res.data.(errx.ValidationErrors)
		if ok {
			for _, message := range errs {
				key := message.Field
				_, ok = resp.Errors[key]
				if !ok {
					resp.Errors[key] = message.Error()
				} else {
					resp.Errors[key] += "; " + message.Error()
				}
			}
		}
	}
	return resp
}

func Invalid(message string, errs errx.ValidationErrors) *InvalidResp {
	if message == "" {
		message = "Ошибка при проверке данных"
	}
	return &InvalidResp{
		Base{false, http.StatusUnprocessableEntity, message, errs},
	}
}

// InternalResp Внутренняя ошибка, связанная с внешними системами HTTPCode = 500.
type InternalResp struct {
	Base
}

func (res *InternalResp) GetHTTPCode() int {
	return http.StatusInternalServerError
}

func (res *InternalResp) GetHTTPResp() interface{} {
	return &struct {
		Status  string      `json:"status"`
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}{
		Status:  "error",
		Code:    res.code,
		Message: "Внутренняя ошибка",
		Data:    res.data,
	}
}

func Internal(message string) *InternalResp {
	return &InternalResp{
		Base{false, http.StatusInternalServerError, message, nil},
	}
}

// FromError перевод ошибки в ответ сервера.
func FromError(err error) Response {
	logErr := errx.Logic{}
	if errors.As(err, &logErr) {
		return BadRequest(logErr.Error(), logErr.Code())
	}
	invErr := errx.Invalid{}
	if errors.As(err, &invErr) {
		return Invalid(invErr.Error(), invErr.Errors())
	}
	nfErr := errx.NotFound{}
	if errors.As(err, &nfErr) {
		return NotFound(invErr.Error(), nfErr.Params)
	}
	base := errx.Base{}
	if errors.As(err, &base) {
		switch base.Kind() {
		case errx.TypePerms:
			return UnAuth()
		case errx.TypeFatal:
			return Internal(base.Error())
		}
	}
	return BadRequest(err.Error(), 400)
}
