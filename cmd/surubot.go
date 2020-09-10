package main

import (
	"log"

	"github.com/annaworks/surubot/pkg/api"
	"github.com/annaworks/surubot/pkg/health"
	"github.com/annaworks/surubot/pkg/slack"
	"github.com/annaworks/surubot/pkg/storage/es"
	Conf "github.com/annaworks/surubot/pkg/conf"

	"go.uber.org/zap"
)

const (
	IndexQuestions = "questions"

	EsQuestionMapping = `{
		"settings": {
			"number_of_shards": 1,
			"number_of_replicas": 1
		},
		"mappings": {
			"properties": {
				"@timestamp": { "type": "date" },
				"question": { "type": "text" },
				"user": { "type": "keyword" }
			}
		}
	}`
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

	// db
	e, err := es.New(c.ES_HOST).Configure()
	if err != nil {
		logger.Fatal("could not create es client", zap.Error(err))
	}

	exists, err := e.DBExists(IndexQuestions)
	if err != nil {
		logger.Fatal("could not check es index", zap.Error(err))
	}
	if !exists {
		err = e.CreateDB(IndexQuestions, EsQuestionMapping)
		if err != nil {
			logger.Fatal("could not create es index", zap.Error(err))
		}
	}

	// api
	api := api.NewApi(logger.Named("api_logger"), c)

	// health service
	h := health.NewHealthService(logger.Named("health_logger"))

	// slack service
	s := slack.NewSlackService(logger.Named("slack_logger"), c)

	api.LoadRoute(h.GetHealthRoute())
	api.LoadRoute(s.GetSlashRoute())
	api.LoadRoute(s.GetInteractionsRoute())
	api.Serve()
}