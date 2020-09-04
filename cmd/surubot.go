package main

import (
	"log"

	"github.com/annaworks/surubot/pkg/api"
	Conf "github.com/annaworks/surubot/pkg/conf"

	"go.uber.org/zap"
)

func main() {
	c := zap.NewProductionConfig()
	c.OutputPaths = []string{"stdout"}
	logger, err := c.Build()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()

	// env
	conf := Conf.NewConf(logger.Named("conf_logger"))

	// api
	api := api.NewApi(logger.Named("api_logger"), conf)
	api.Init()
	api.Serve()
}