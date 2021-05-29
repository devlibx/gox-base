package goxServer

import (
	"github.com/devlibx/gox-base"
	"github.com/devlibx/gox-base/config"
	"net/http"
)

type Server interface {
	Start(handler http.Handler, appConfig *config.App) error
}

func NewServer(cf gox.CrossFunction) (Server, error) {
	s := &serverImpl{CrossFunction: cf}
	return s, nil
}
