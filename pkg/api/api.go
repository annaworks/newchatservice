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
	Logger *zap.Logger
	Router *mux.Router
	Conf	Conf.Conf
}

func NewApi(logger *zap.Logger, conf Conf.Conf) *Api {
	return &Api{
		Logger: logger,
		Router: mux.NewRouter().StrictSlash(true),
		Conf: conf,
	}
}

type Route struct {
	Path string
	Handler http.HandlerFunc
	Method string
	Name string
}

func (a *Api) LoadRoute(r *Route) *Api {
	a.Router.HandleFunc(r.Path, r.Handler).
		Methods(r.Method).
		Name(r.Name)
	
	a.Logger.Info(fmt.Sprintf("route loaded: %s", r.Name))
	return a
}

func (a *Api) Serve() {
	s := http.Server{
		Addr:           fmt.Sprintf(":%s", a.Conf.API_PORT),
		Handler:        a.Router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   20 * time.Second,
		IdleTimeout: 	120 * time.Second,
	}
	
	a.Logger.Info(fmt.Sprintf("Server serving on port %s", a.Conf.API_PORT))

	if err := s.ListenAndServe(); err != nil {
		a.Logger.Fatal("Error in serving server", zap.Error(err))
	}
}