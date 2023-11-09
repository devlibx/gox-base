package goxServer

import (
	"context"
	"fmt"
	"github.com/devlibx/gox-base"
	"github.com/devlibx/gox-base/config"
	"go.uber.org/fx"
	"net/http"
	"sync"
)

type Server interface {
	Start(handler http.Handler, appConfig *config.App) error
	Stop() chan bool
}

type ServerShutdownHook interface {
	Setup(interface{})
	StopFunction() func()
}

func NewServer(cf gox.CrossFunction) (Server, error) {
	s := &serverImpl{CrossFunction: cf, stopOnce: &sync.Once{}}
	return s, nil
}

func NewServerWithShutdownHookFunc(cf gox.CrossFunction, shutdownHookFunc func()) (Server, error) {
	s := &serverImpl{CrossFunction: cf, shutdownHookFunc: shutdownHookFunc, stopOnce: &sync.Once{}}
	return s, nil
}

func NewServerWithShutdownHook(cf gox.CrossFunction, serverShutdownHook ServerShutdownHook) (Server, error) {
	s := &serverImpl{CrossFunction: cf, shutdownHookFunc: serverShutdownHook.StopFunction(), stopOnce: &sync.Once{}}
	return s, nil
}

// noOpServerShutdownHook - dummy shutdown hook
type noOpServerShutdownHook struct {
}

func (h noOpServerShutdownHook) Setup(set interface{}) {
}

func (noOpServerShutdownHook) StopFunction() func() {
	return func() {
	}
}

func NoOpServerShutdownHook() ServerShutdownHook {
	return &noOpServerShutdownHook{}
}

// FxServerShutdownHook is a common fx
type FxServerShutdownHook struct {
	FxApp *fx.App
}

func (h *FxServerShutdownHook) Setup(i interface{}) {
	if a, ok := i.(*fx.App); ok {
		h.FxApp = a
	} else {
		fmt.Println("*********** Something is wrong - we expected *fx.App in FxServerShutdownHook.Setup() *********** ")
	}
}

func (h *FxServerShutdownHook) StopFunction() func() {
	return func() {
		if h.FxApp != nil {
			_ = h.FxApp.Stop(context.Background())
		}
	}
}

func NewFxServerShutdownHook() ServerShutdownHook {
	return &FxServerShutdownHook{}
}
