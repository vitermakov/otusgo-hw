package rest

import (
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/utils/errx"
	"net/http"
)

// Response интерфейс для связи с веб-сервером
type Response interface {
	GetHttpCode() int         // HTTP-код ответа
	GetHttpResp() interface{} // содержимое ответа
	Message() string
	Success() bool
}

// response общий вид ответа сервера
type response struct {
	success bool
	code    int
	message string
	data    interface{}
}

func (res *response) GetHttpCode() int {
	if res.success {
		return http.StatusOK
	}
	return http.StatusBadRequest
}
func (res *response) GetHttpResp() interface{} {
	var resp = &struct {
		Status  string      `json:"status"`
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
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
func (res response) Success() bool {
	return res.success
}
func (res response) Message() string {
	return res.message
}

// respOK Удачный ответ на запросы (команды) типа POST/PUT/PATCH/DELETE HTTPCode = 200
type respOK struct {
	response
}

func (res respOK) GetHttpCode() int {
	return http.StatusOK
}
func OK(message string, data interface{}) *respOK {
	if len(message) <= 0 {
		message = "Успешно"
	}
	return &respOK{response{true, http.StatusOK, message, data}}
}

// respData Удачный ответ на запрос GET на получение объекта HTTPCode = 200
type respData struct {
	data interface{}
}

func (res *respData) GetHttpCode() int {
	return http.StatusOK
}
func (res respData) Success() bool {
	return true
}
func (res respData) Message() string {
	return ""
}
func (res respData) GetHttpResp() interface{} {
	return res.data
}
func Data(data interface{}) *respData {
	return &respData{data}
}

// respList Удачный ответ на запрос GET на получение списка
// объектов без навигации HTTPCode = 200
type respList struct {
	items interface{}
}

func (res *respList) GetHttpCode() int {
	return http.StatusOK
}
func (res respList) Success() bool {
	return true
}
func (res respList) Message() string {
	return ""
}
func (res respList) GetHttpResp() interface{} {
	return res.items
}
func List(list interface{}) *respList {
	return &respList{list}
}

// respBad Ошибка из-за нарушения правил бизнес-логики HTTPCode = 400
// Ошибки, связанные с действиями пользователей, которые не могут быть выполнены при текущих правилах
type respBad struct {
	response
}

func BadRequest(message string, code int) *respBad {
	return &respBad{response{false, code, message, nil}}
}

// respUnAuth Попытка совершить операцию без авторизации HTTPCode = 401
type respUnAuth struct {
	response
}

func (res *respUnAuth) GetHttpCode() int {
	return http.StatusUnauthorized
}
func UnAuth() *respUnAuth {
	message := "Вы не авторизованы"
	return &respUnAuth{
		response{false, http.StatusUnauthorized, message, nil},
	}
}

// respNotFound Запрошенный ресурс не найден HTTPCode = 404
type respNotFound struct {
	response
}

func (res *respNotFound) GetHttpCode() int {
	return http.StatusNotFound
}
func NotFound(message string, params interface{}) *respNotFound {
	if len(message) <= 0 {
		message = "Объект не найден"
	}
	return &respNotFound{
		response{false, http.StatusNotFound, message, params},
	}
}

// respInvalid Ошибка валидации входных данных HTTPCode = 422
type respInvalid struct {
	response
}

func (res *respInvalid) GetHttpCode() int {
	return http.StatusUnprocessableEntity
}
func (res respInvalid) GetHttpResp() interface{} {
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
		errors, ok := res.data.([]errx.ValidationError)
		if ok {
			for _, message := range errors {
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
func Invalid(message string, errs errx.ValidationErrors) *respInvalid {
	if len(message) <= 0 {
		message = "Ошибка при проверке данных"
	}
	return &respInvalid{
		response{false, http.StatusUnprocessableEntity, message, errs},
	}
}

// respInternal Внутренняя ошибка, связанная с внешними системами HTTPCode = 500
type respInternal struct {
	response
}

func (res *respInternal) GetHttpCode() int {
	return http.StatusInternalServerError
}
func (res *respInternal) GetHttpResp() interface{} {
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
func Internal(message string) *respInternal {
	return &respInternal{
		response{false, http.StatusInternalServerError, message, nil},
	}
}

// NewErrorResponse перевод ошибки в ответ сервера
func NewErrorResponse(err error) Response {
	switch err.(type) {
	case errx.Logic:
		return BadRequest(err.Error(), err.(errx.Logic).Code())
	case errx.Invalid:
		return Invalid(err.Error(), err.(errx.Invalid).Errors())
	case errx.Base:
		bsErr := err.(errx.Base)
		switch bsErr.Kind() {
		case errx.TypePerms:
			return UnAuth()
		case errx.TypeFatal:
			return Internal(bsErr.Error())
		}
	}
	return BadRequest(err.Error(), 400)
}
