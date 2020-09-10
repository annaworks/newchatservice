package slack

import (
	"go.uber.org/zap"
	Conf "github.com/annaworks/surubot/pkg/conf"
	"github.com/slack-go/slack"
)

type SlackService struct {
	Logger *zap.Logger
	Conf Conf.Conf
	Api *slack.Client
}

func NewSlackService(logger *zap.Logger, conf Conf.Conf) *SlackService {
	return &SlackService{
		Logger: logger,
		Conf: conf,
		Api: slack.New(conf.SLACK_TOKEN),
	}
}