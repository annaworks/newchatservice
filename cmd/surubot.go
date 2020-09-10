package main

import (
	"log"

	"github.com/annaworks/surubot/pkg/api"
	"github.com/annaworks/surubot/pkg/health"
	"github.com/annaworks/surubot/pkg/slack"
	"github.com/annaworks/surubot/pkg/storage/es"
	"github.com/annaworks/surubot/pkg/knowledge"
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

	// db
	e, err := es.New(c.ES_HOST).Configure()
	if err != nil {
		logger.Fatal("could not create es client", zap.Error(err))
	}

	exists, err := e.IndexExists(knowledge.IndexQuestions)
	if err != nil {
		logger.Fatal("could not check es index", zap.Error(err))
	}
	if !exists {
		err = e.CreateIndex(knowledge.IndexQuestions, knowledge.EsQuestionMapping)
		if err != nil {
			logger.Fatal("could not create es index", zap.Error(err))
		}
	}

	exists, err = e.IndexExists(knowledge.IndexAnswers)
	if err != nil {
		logger.Fatal("could not check es index", zap.Error(err))
	}
	if !exists {
		err = e.CreateIndex(knowledge.IndexAnswers, knowledge.EsAnswerMapping)
		if err != nil {
			logger.Fatal("could not create es index", zap.Error(err))
		}
	}

	// api
	api := api.NewApi(logger.Named("api_logger"), c)

	// health service
	h := health.NewHealthService(logger.Named("health_logger"))

	// slack service
	s := slack.NewSlackService(logger.Named("slack_logger"), c, e)

	api.LoadRoute(h.GetHealthRoute())
	api.LoadRoute(s.GetSlashRoute())
	api.LoadRoute(s.GetInteractionsRoute())
	api.Serve()
}