package knowledge

import (
	"time"
)

const (
	IndexAnswers = "answers"

	EsAnswerMapping = `{
		"settings": {
			"number_of_shards": 1,
			"number_of_replicas": 1
		},
		"mappings": {
			"properties": {
				"@timestamp": { "type": "date" },
				"question_id": { "type": "keyword" },
				"question": { "type": "text" },
				"answer": { "type": "text" },
				"user": { "type": "keyword" }
			}
		}
	}`
)

type Answer struct {
	Timestamp  string `json:"@timestamp"`
	Question   string `json:"question"`
	QuestionId string `json:"question_id"`
	Answer     string `json:"answer"`
	User       string `json:"user"`
}

func NewAnswer(question, questionId, answer, user string) *Answer {
	return &Answer{
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
		Question:   question,
		QuestionId: questionId,
		Answer:     answer,
		User:       user,
	}
}
