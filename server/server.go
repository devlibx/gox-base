package goxServer

import (
	"github.com/devlibx/gox-base"
	"github.com/devlibx/gox-base/config"
	"net/http"
	"sync"
)

type Server interface {
	Start(handler http.Handler, appConfig *config.App) error
	Stop() chan bool
}

func NewServer(cf gox.CrossFunction) (Server, error) {
	s := &serverImpl{CrossFunction: cf, stopOnce: &sync.Once{}}
	return s, nil
}
