package knowledge

import (
	"time"
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

type Question struct {
	Timestamp string `json:"@timestamp"`
	Question  string `json:"question"`
	User      string `json:"user"`
}

func NewQuestion(question, user string) *Question {
	return &Question{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Question: question,
		User: user,
	}
}