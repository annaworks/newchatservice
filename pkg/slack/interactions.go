package slack

import (
	"net/http"
	"fmt"
	"encoding/json"
	"errors"

	"github.com/annaworks/surubot/pkg/api"
	"github.com/slack-go/slack"
	"go.uber.org/zap"
)

func (s SlackService) GetInteractionsRoute() *api.Route{
	return &api.Route{
		Path: "/api/v1/interactions",
		Handler: s.HandleInteractionsRequest,
		Method: http.MethodPost,
		Name: "interactions",
	}
}

func (s SlackService) processPayload(r *http.Request) (*ButtonActionPayload, error) {
	var payload slack.InteractionCallback
	err := json.Unmarshal([]byte(r.FormValue("payload")), &payload)
	if err != nil {
		return &ButtonActionPayload{}, fmt.Errorf("Could not parse ineraction callback payload: %v", err)
	}

	message, err := s.getQuestionText(payload)
	if err != nil {
		return &ButtonActionPayload{}, fmt.Errorf("Could not find question text from interaction callback payload: %v", err)
	}

	return newButtonActionPayload(
		message,
		payload.User.Name, 
		payload.User.ID,
		payload.ActionCallback.BlockActions[0].Value,
		payload.TriggerID,
		payload.Type,
	), nil
}

func (s SlackService) HandleInteractionsRequest(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info("Received an interactions request")

	p, err := s.processPayload(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Error: Unknown action")))
		s.Logger.Error("Error processing interaction callback payload", zap.Error(err))
		return
	}
	fmt.Printf("Payload: %+v", p)
}

type ButtonActionPayload struct {
	Question string
	Username string
	UserID string
	ButtonClicked string
	TriggerID string
	Type slack.InteractionType
}

func newButtonActionPayload(question, username, userID, buttonClicked, triggerID string, interactionType slack.InteractionType) *ButtonActionPayload {
	return &ButtonActionPayload {
		Question: question,
		Username: username,
		UserID: userID,
		ButtonClicked: buttonClicked,
		TriggerID: triggerID,
		Type: interactionType,
	}
}

func (s SlackService) getQuestionText(payload slack.InteractionCallback) (string, error)  {
	var message string
	for _, b := range payload.Message.Msg.Blocks.BlockSet {
		if b.BlockType() == "section"  {
			s := b.(*slack.SectionBlock)
			message = s.Fields[0].Text
			break
		}
	}
	if message == "" {
		return message, errors.New(fmt.Sprintf("Error: Unknown block type found in block action"))
	}
	return message, nil
}