package slack

import (
	"go.uber.org/zap"
	"github.com/annaworks/surubot/pkg/storage/es"
	Conf "github.com/annaworks/surubot/pkg/conf"
	"github.com/slack-go/slack"
)

type SlackService struct {
	Logger *zap.Logger
	Conf Conf.Conf
	Api *slack.Client
	Storage *es.ES
}

func NewSlackService(logger *zap.Logger, conf Conf.Conf, e *es.ES) *SlackService {
	return &SlackService{
		Logger: logger,
		Conf: conf,
		Api: slack.New(conf.SLACK_TOKEN),
		Storage: e,
	}
}