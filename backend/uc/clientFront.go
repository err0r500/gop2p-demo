package uc

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"gop2p/domain"
	"log"
)

// ClientFrontLogic handles the logic exposed to the frontend
type ClientFrontLogic interface {
	NewSessionRegistered(username string) error
	SendMessageToOtherClient(toUserName string, msg string) error
	GetConversationWith(authorName string) ([]domain.Message, error)
}

type clientFrontInteractor struct {
	currentUsername string
	cm              ConversationManager
	sg              ServerGateway
	cg              ClientGateway
}

func NewClientFrontLogic(cm ConversationManager, sg ServerGateway, cg ClientGateway) ClientFrontLogic {
	return &clientFrontInteractor{
		currentUsername: "",
		cm:              cm,
		sg:              sg,
		cg:              cg,
	}
}

// RegisterNewSession is used by the client to register a new session
func (i *clientFrontInteractor) NewSessionRegistered(username string) error {
	tracer := opentracing.GlobalTracer()
	span := tracer.StartSpan("http:new_session_ok")
	defer span.Finish()

	log.Println("UC : newSesionRegistered", username)
	i.currentUsername = username
	return nil
}

// SendMessageToOtherClient is used by the client to send a message to another one
func (i clientFrontInteractor) SendMessageToOtherClient(toUserName string, msg string) error {
	span := opentracing.GlobalTracer().StartSpan("uc:send_message_to_other_client")
	defer span.Finish()

	emitter := i.currentUsername
	if emitter == "" {
		log.Println("unauthorized currentUser is empty")
		return domain.ErrUnauthorized{}
	}

	s, err := i.sg.AskSessionToServer(opentracing.ContextWithSpan(context.Background(), span), emitter, toUserName)
	if err != nil {
		return domain.ErrTechnical{}
	}
	if s == nil {
		return domain.ErrResourceNotFound{}
	}

	if err := i.cm.AppendToConversationWith(toUserName, emitter, msg); err != nil {
		return domain.ErrTechnical{}
	}

	if err := i.cg.SendMsg(s.Address,
		domain.Message{Author: emitter, Content: msg},
		i.currentUsername,
	); err != nil {
		return domain.ErrTechnical{}
	}

	log.Println("sent new message to", toUserName)
	return nil
}

// GetConversationWith is used by the client to get a given conversation
func (i clientFrontInteractor) GetConversationWith(authorName string) ([]domain.Message, error) {
	span := opentracing.GlobalTracer().StartSpan("uc:get_conversation_with")
	defer span.Finish()

	messages, err := i.cm.GetConversationWith(authorName)
	if err != nil {
		return nil, domain.ErrTechnical{}
	}

	return messages, nil
}
