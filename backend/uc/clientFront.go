package uc

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"gop2p/domain"
)

// ClientFrontLogic handles the logic exposed to the frontend
type ClientFrontLogic interface {
	NewSessionRegistered(ctx context.Context, username string) error
	SendMessageToOtherClient(ctx context.Context, toUserName string, msg string) error
	GetConversationWith(ctx context.Context, authorName string) ([]domain.Message, error)
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
func (i *clientFrontInteractor) NewSessionRegistered(ctx context.Context, username string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "uc:new_session_registered")
	defer span.Finish()

	i.currentUsername = username
	return nil
}

// SendMessageToOtherClient is used by the client to send a message to another one
func (i clientFrontInteractor) SendMessageToOtherClient(ctx context.Context, toUserName string, msg string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "uc:send_message_to_other_client")
	defer span.Finish()

	emitter := i.currentUsername
	if emitter == "" {
		span.LogFields(log.Error(errors.New("missing current user session")))
		return domain.ErrUnauthorized{}
	}

	s, ok := i.sg.AskSessionToServer(opentracing.ContextWithSpan(context.Background(), span), emitter, toUserName)
	if !ok {
		return domain.ErrTechnical{}
	}
	if s == nil {
		span.LogFields(log.Error(errors.New("session not found")))
		return domain.ErrResourceNotFound{}
	}

	if ok := i.cm.AppendToConversationWith(ctx, toUserName, emitter, msg); !ok {
		return domain.ErrTechnical{}
	}

	if ok := i.cg.SendMsg(ctx, s.Address,
		domain.Message{Author: emitter, Content: msg},
		i.currentUsername,
	); !ok {
		return domain.ErrTechnical{}
	}

	return nil
}

// GetConversationWith is used by the client to get a given conversation
func (i clientFrontInteractor) GetConversationWith(ctx context.Context, authorName string) ([]domain.Message, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "uc:get_conversation_with")
	defer span.Finish()

	messages, ok := i.cm.GetConversationWith(ctx, authorName)
	if !ok {
		return nil, domain.ErrTechnical{}
	}

	return messages, nil
}
