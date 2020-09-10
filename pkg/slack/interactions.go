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

func (s SlackService) HandleInteractionsRequest(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info("Received an interactions request")

	p, err := processPayload(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Error: Unknown action")))
		s.Logger.Error("Error processing interaction callback payload", zap.Error(err))
		return
	}
	fmt.Printf("Payload: %+v", p)

	switch p.Type {
	case slack.InteractionTypeBlockActions:
		switch p.ButtonClicked {
		case viewButtonID:
			modalRequest := p.newViewRequest()
			_, err = s.Api.OpenView(p.TriggerID, modalRequest)
			if err != nil {
				fmt.Printf("Error opening view: %s", err)
			}
		case answerButtonID:
			modalRequest := p.newAnswerRequest()
			_, err = s.Api.OpenView(p.TriggerID, modalRequest)
			if err != nil {
				fmt.Printf("Error opening view: %s", err)
			}
		default:
			fmt.Printf("Unknown interaction type block action: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Unknown interaction type block action")))
			s.Logger.Error("Error receiving unknown interaction type block action")
		}
	default:
		fmt.Printf("Unknown interaction callback type: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Error: Unknown interaction callback type")))
		s.Logger.Error("Error receiving unknown interaction callback type")
	}
}

func processPayload(r *http.Request) (*ButtonActionPayload, error) {
	var payload slack.InteractionCallback
	err := json.Unmarshal([]byte(r.FormValue("payload")), &payload)
	if err != nil {
		return &ButtonActionPayload{}, fmt.Errorf("Could not parse ineraction callback payload: %v", err)
	}

	message, err := getQuestionText(payload)
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

func getQuestionText(payload slack.InteractionCallback) (string, error)  {
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

func (p ButtonActionPayload) newViewRequest() slack.ModalViewRequest {
	titleText := slack.NewTextBlockObject("plain_text", "View answers", false, false)
	closeText := slack.NewTextBlockObject("plain_text", "Close", false, false)
	submitText := slack.NewTextBlockObject("plain_text", "Add answer", false, false)

	headerContent := fmt.Sprintf("*Question:*\n %v", p.Question)
	headerText := slack.NewTextBlockObject("mrkdwn", headerContent, false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	// Empty state on initial view
	emptyContentText := slack.NewTextBlockObject("plain_text", "There are currently no answers for this question.", false, false)
	emptyContentSection := slack.NewSectionBlock(emptyContentText, nil, nil)

	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			headerSection,
			emptyContentSection,
		},
	}

	var modalRequest slack.ModalViewRequest
	modalRequest.Type = slack.ViewType("modal")
	modalRequest.Title = titleText
	modalRequest.Close = closeText
	modalRequest.Submit = submitText
	modalRequest.Blocks = blocks
	return modalRequest
}

func (p ButtonActionPayload) newAnswerRequest() slack.ModalViewRequest {
	titleText := slack.NewTextBlockObject("plain_text", "Add an answer", false, false)
	closeText := slack.NewTextBlockObject("plain_text", "Cancel", false, false)
	submitText := slack.NewTextBlockObject("plain_text", "Submit", false, false)

	answerText := slack.NewTextBlockObject("plain_text", p.Question, false, false)
	answerPlaceholder := slack.NewTextBlockObject("plain_text", "Write something", false, false)
	answerElement := slack.NewPlainTextInputBlockElement(answerPlaceholder, "answer")
	answerElement.Multiline = true
	answer := slack.NewInputBlock("Answer", answerText, answerElement)

	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			answer,
		},
	}

	var modalRequest slack.ModalViewRequest
	modalRequest.Type = slack.ViewType("modal")
	modalRequest.Title = titleText
	modalRequest.Close = closeText
	modalRequest.Submit = submitText
	modalRequest.Blocks = blocks
	return modalRequest
}