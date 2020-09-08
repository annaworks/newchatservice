package main

import (
	"log"

	"github.com/annaworks/surubot/pkg/api"
	"github.com/annaworks/surubot/pkg/health"
	"github.com/annaworks/surubot/pkg/slack"
	Conf "github.com/annaworks/surubot/pkg/conf"

	"go.uber.org/zap"
)

func main() {
	z := zap.NewProductionConfig()
	z.OutputPaths = []string{"stdout"}
	logger, err := z.Build()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()

	// env
	c := Conf.NewConf(logger.Named("conf_logger"))

	// api
	api := api.NewApi(logger.Named("api_logger"), c)

	// health service
	h := health.NewHealthService(logger.Named("health_logger"))

	// slack service
	s := slack.NewSlashService(logger.Named("slack_slash_logger"), c)

	api.LoadRoute(h.GetHealthRoute())
	api.LoadRoute(s.GetSlashRoute())
	api.Serve()
}