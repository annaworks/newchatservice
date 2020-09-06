package api

import (
	"fmt"
	"net/http"
	"time"

	Conf "github.com/annaworks/surubot/pkg/conf"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Api struct {
	logger *zap.Logger
	router *mux.Router
	conf	Conf.Conf
}

func NewApi(logger *zap.Logger, conf Conf.Conf) *Api {
	return &Api{
		logger: logger,
		router: mux.NewRouter().StrictSlash(true),
		conf: conf,
	}
}

type Route struct {
	Path string
	Handler http.HandlerFunc
	Method string
	Name string
}

func (a *Api) LoadRoute(r *Route) *Api {
	a.router.HandleFunc(r.Path, r.Handler).
		Methods(r.Method).
		Name(r.Name)
	
	a.logger.Info(fmt.Sprintf("route loaded: %s", r.Name))
	return a
}

func (a *Api) Serve() {
	s := http.Server{
		Addr:           fmt.Sprintf(":%s", a.conf.API_PORT),
		Handler:        a.router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   20 * time.Second,
		IdleTimeout: 	120 * time.Second,
	}
	
	a.logger.Info(fmt.Sprintf("Server serving on port %s", a.conf.API_PORT))

	if err := s.ListenAndServe(); err != nil {
		a.logger.Fatal("Error in serving server", zap.Error(err))
	}
}