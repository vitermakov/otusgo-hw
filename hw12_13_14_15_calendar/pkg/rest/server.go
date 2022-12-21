package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/logger"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/utils/errx"
	"log"
	"net"
	"net/http"
	"runtime/debug"
	"strconv"
	"time"
)

const (
	DefaultHost = "localhost"
	DefaultPort = 8080
)

type Config struct {
	host  string
	port  int
	debug bool
}

func (cfg Config) Host() string {
	if len(cfg.host) > 0 {
		return cfg.host
	}
	return DefaultHost
}

func (cfg Config) Port() int {
	if cfg.port > 0 {
		return cfg.port
	}
	return DefaultPort
}

func (cfg Config) Debug() bool {
	return cfg.debug
}

// Server простой Http-сервер для того, чтобы гонять Json через Http.
type Server struct {
	Logger      logger.Logger
	AuthService AuthService
	config      Config
	engine      *http.Server
	router      *mux.Router
}

type HandlerFunc func(r *http.Request) Response

func NewServer(cfg Config, authSrv AuthService, logger logger.Logger) *Server {
	return &Server{
		Logger:      logger,
		AuthService: authSrv,
		router:      mux.NewRouter(),
		config:      cfg,
	}
}

func (s *Server) Start() error {
	s.router.Use(
		s.stopPanic,
		s.loggingMiddleware,
		s.authMiddleware,
	)
	listenAddress := net.JoinHostPort(s.config.Host(), strconv.Itoa(s.config.Port()))
	s.engine = &http.Server{
		Addr:    listenAddress,
		Handler: s.router,
	}
	err := s.engine.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	err := s.engine.Shutdown(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) GET(pattern string, handler HandlerFunc) {
	s.router.HandleFunc(pattern, s.WrapHandler(handler)).Methods("GET")
}

func (s *Server) POST(pattern string, handler HandlerFunc) {
	s.router.HandleFunc(pattern, s.WrapHandler(handler)).Methods("POST")
}

func (s *Server) PUT(pattern string, handler HandlerFunc) {
	s.router.HandleFunc(pattern, s.WrapHandler(handler)).Methods("PUT")
}

func (s *Server) DELETE(pattern string, handler HandlerFunc) {
	s.router.HandleFunc(pattern, s.WrapHandler(handler)).Methods("DELETE")
}

func (s *Server) WrapHandler(handler HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var response Response
		doneChan := make(chan bool)
		go func() {
			response = handler(request)
			close(doneChan)
		}()
		select {
		case <-request.Context().Done():
			response = NewErrorResponse(errx.FatalNew("Abort in timeout"))
			log.Println("Abort in timeout", request.Method, request.URL)
		case <-doneChan:
		}
		s.showResponse(writer, response)
	}
}

func (s *Server) showResponse(w http.ResponseWriter, r Response) {
	w.WriteHeader(r.GetHttpCode())
	data, _ := json.Marshal(r.GetHttpResp())
	_, _ = w.Write(data)
}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timeStart := time.Now()
		next.ServeHTTP(w, r)
		s.Logger.Info(
			map[string]interface{}{
				"method":  r.Method,
				"url":     r.URL,
				"latency": time.Now().Sub(timeStart).Milliseconds(),
			},
			"http query",
		)
	})
}

func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		authCredentials := r.Header.Get("Authorization")
		user, err := s.AuthService.Authorize(ctx, authCredentials)
		if err != nil {
			panic(fmt.Sprintf("error auth-service: %v", err))
		}
		if user != nil {
			r = r.WithContext(context.WithValue(ctx, "user", map[string]string{
				"id":    user.ID,
				"name":  user.Name,
				"login": user.Login,
			}))
			next.ServeHTTP(w, r)
			return
		}

		response := NewErrorResponse(errx.PermsNew("нет доступа"))
		s.showResponse(w, response)
	})
}

func (s *Server) stopPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				response := NewErrorResponse(errx.FatalNew(fmt.Sprintf("%+v\n%+v", r, string(debug.Stack()))))
				s.showResponse(w, response)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
