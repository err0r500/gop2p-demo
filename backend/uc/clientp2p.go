package uc

import (
	"github.com/opentracing/opentracing-go"
	"gop2p/domain"
	"log"
)

// ClientP2PLogic handles the logic of the central server
type ClientP2PLogic interface {
	HandleMessageReceived(msg string, emitter domain.User) error
}

type clientp2pInteractor struct {
	cm ConversationManager
}

func NewClientP2pLogic(cm ConversationManager) ClientP2PLogic {
	return clientp2pInteractor{cm: cm}
}

// HandleMessageReceived is used by the client to handle a new message
func (i clientp2pInteractor) HandleMessageReceived(msg string, emitter domain.User) error {
	span := opentracing.GlobalTracer().StartSpan("uc:handle_new_message_received")
	defer span.Finish()

	if err := i.cm.AppendToConversationWith(emitter.Login, emitter.Login, msg); err != nil {
		log.Println(err)
		return domain.ErrTechnical{}
	}

	return nil
}
