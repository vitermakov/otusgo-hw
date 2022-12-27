package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/felixge/httpsnoop"
	"github.com/gorilla/mux"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/logger"
	rs "github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/rest/rqres"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/utils/errx"
)

const (
	DefaultHost = "localhost"
	DefaultPort = 8080
)

type CtxKey struct{}

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

type HandlerFunc func(r *rs.Request) rs.Response

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
		s.loggingMiddleware,
		s.authMiddleware,
	)
	listenAddress := net.JoinHostPort(s.config.Host(), strconv.Itoa(s.config.Port()))
	s.engine = &http.Server{
		Addr:              listenAddress,
		Handler:           s.router,
		ReadTimeout:       1 * time.Second,
		WriteTimeout:      1 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
	}
	err := s.engine.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if err := s.engine.Shutdown(ctx); err != nil {
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
	return func(writer http.ResponseWriter, rq *http.Request) {
		var response rs.Response
		if rq.Method == "OPTIONS" {
			writer.WriteHeader(http.StatusOK)
			return
		}
		ctx, cancel := context.WithTimeout(rq.Context(), 30*time.Second)
		defer cancel()
		rq = rq.WithContext(ctx)

		request := &rs.Request{Request: rq, Params: mux.Vars(rq)}
		doneChan := make(chan bool)
		go func() {
			defer func() {
				if r := recover(); r != nil {
					response = rs.FromError(errx.FatalNew(fmt.Errorf("%+v\n%+v", r, string(debug.Stack()))))
					close(doneChan)
				}
			}()
			response = handler(request)
			close(doneChan)
		}()
		select {
		case <-request.Context().Done():
			err := errors.New("abort in timeout")
			response = rs.FromError(errx.FatalNew(err))
			s.Logger.Error("Abort in timeout (%ds). %s %s, ", s.engine.IdleTimeout.Milliseconds(), request.Method, request.URL)
		case <-doneChan:
		}
		s.showResponse(writer, response)
	}
}

func (s *Server) showResponse(w http.ResponseWriter, resp rs.Response) {
	w.Header().Set("Content-type", "application/json; charset=utf-8")
	w.WriteHeader(resp.GetHTTPCode())
	data, _ := json.Marshal(resp.GetHTTPResp())
	_, _ = w.Write(data)
}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var statusCode int
		timeStart := time.Now()
		wrapped := httpsnoop.Wrap(w, httpsnoop.Hooks{
			WriteHeader: func(httpsnoop.WriteHeaderFunc) httpsnoop.WriteHeaderFunc {
				return func(code int) {
					w.WriteHeader(code)
					statusCode = code
				}
			},
		})
		next.ServeHTTP(wrapped, req)
		s.Logger.Info(
			fmt.Sprintf(
				"%s %s %s %s %d %s \"%s\"",
				strings.Split(req.RemoteAddr, ":")[0],
				req.Method,
				req.RequestURI,
				req.Proto,
				statusCode,
				time.Since(timeStart).String(),
				req.Header.Get("User-Agent"),
			),
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
			r = r.WithContext(context.WithValue(ctx, CtxKey{}, map[string]string{ //nolint:go-staticcheck // fdddd
				"id":    user.ID,
				"name":  user.Name,
				"login": user.Login,
			}))
			next.ServeHTTP(w, r)
			return
		}

		response := rs.FromError(errx.PermsNew(fmt.Errorf("нет доступа")))
		s.showResponse(w, response)
	})
}
