package rqres

import "net/http"

// Request расширение стандартного http-запроса
type Request struct {
	*http.Request
	Params map[string]string // именованные параметры из path
}

func (rq Request) Param(key string) string {
	return rq.Params[key]
}
