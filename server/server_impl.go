package goxServer

import (
	"context"
	"fmt"
	"net/http"
	"net/http/pprof"
	"strings"
	"sync"
	"time"

	"github.com/devlibx/gox-base"
	"github.com/devlibx/gox-base/config"
	"github.com/devlibx/gox-base/errors"
	"github.com/urfave/negroni"
	"go.uber.org/zap"
	"gopkg.in/tylerb/graceful.v1"
)

type serverImpl struct {
	server           *http.Server
	gracefulServer   *graceful.Server
	serverRunning    chan bool
	stopOnce         *sync.Once
	shutdownHookFunc func()
	gox.CrossFunction
}

func (s *serverImpl) Start(handler http.Handler, applicationConfig *config.App) error {
	if applicationConfig == nil {
		return errors.New("application config is nil")
	}

	// Channel to wait for server to stop
	s.serverRunning = make(chan bool, 1)

	// Setup default values
	applicationConfig.SetupDefaults()

	// Setup server
	var rootHandler *negroni.Negroni
	if applicationConfig.EnablePProf {
		rootHandler = negroni.New(negroni.NewRecovery(), negroni.NewStatic(http.Dir("public")), newPprofHandler())
	} else {
		rootHandler = negroni.New(negroni.NewRecovery(), negroni.NewStatic(http.Dir("public")))
	}
	if applicationConfig.IsServerTimeLoggingEnabled() {
		rootHandler.Use(s.setupTimeLogging())
	}
	rootHandler.UseHandler(handler)

	// Setup http server
	s.server = &http.Server{
		Handler:      rootHandler,
		Addr:         fmt.Sprintf("0.0.0.0:%d", applicationConfig.HttpPort),
		WriteTimeout: time.Duration(applicationConfig.RequestWriteTimeoutMs) * time.Millisecond,
		ReadTimeout:  time.Duration(applicationConfig.RequestReadTimeoutMs) * time.Millisecond,
		IdleTimeout:  time.Duration(applicationConfig.IdleTimeoutMs) * time.Millisecond,
	}

	s.gracefulServer = &graceful.Server{
		Timeout:           time.Duration(applicationConfig.OutstandingRequestTimeoutMs) * time.Second,
		Server:            s.server,
		ShutdownInitiated: s.shutdownHookFunc,
		NoSignalHandling: true,
	}

	return s.gracefulServer.ListenAndServe()
}

func (s *serverImpl) Stop() chan bool {
	s.stopOnce.Do(func() {
		go func() {
			_ = s.gracefulServer.Shutdown(context.TODO())
			s.serverRunning <- true
			close(s.serverRunning)
		}()
	})
	return s.serverRunning
}

func (s *serverImpl) setupTimeLogging() negroni.HandlerFunc {
	logger := s.Logger().Named("negroni")
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		start := time.Now()
		next(rw, r)
		end := time.Now()
		logger.Info("",
			zap.String("remoteAddr", r.RemoteAddr),
			zap.String("source", r.Header.Get("X-FORWARDED-FOR")),
			zap.Int64("duration", end.Sub(start).Milliseconds()),
		)
	}
}

type pprofHandler struct {
	mux *http.ServeMux
}

func (rec *pprofHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if strings.Contains(r.RequestURI, "/debug/pprof") {
		rec.mux.ServeHTTP(rw, r)
	} else {
		next(rw, r)
	}
}

func newPprofHandler() *pprofHandler {
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	return &pprofHandler{mux: mux}
}
