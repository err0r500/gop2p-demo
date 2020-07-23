package conversationmanager

import (
	"context"
	"errors"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"gop2p/domain"
	"gop2p/uc"
	"sync"
)

type store struct {
	rw            *sync.Map
	failingMethod string
}

// New is the constructor of this in memory implementation of the uc.SessionManager
func New() uc.ConversationManager {
	return store{rw: &sync.Map{}}
}

type FailingConversationManager interface {
	uc.ConversationManager
	InjectErrorAt(failingMethod string)
}

func NewFailable() FailingConversationManager {
	return &store{rw: &sync.Map{}, failingMethod: ""}
}

func (s *store) InjectErrorAt(failingMethod string) {
	s.failingMethod = failingMethod
}

func (s store) GetConversationWith(ctx context.Context, authorName string) ([]domain.Message, bool) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "conversation_manager:get-conversation_with")
	defer span.Finish()

	if s.failingMethod == "getConversationWith" {
		return nil, false
	}

	val, ok := s.rw.Load(authorName)
	if !ok {
		return nil, true
	}

	conversation, ok := val.([]domain.Message)
	if !ok {
		span.LogFields(log.Error(errors.New("not a conversation stored at Key")))
		return nil, false
	}
	return conversation, true
}

func (s store) AppendToConversationWith(ctx context.Context, userName, msgAuthor, msgContent string) bool {
	span, ctx := opentracing.StartSpanFromContext(ctx, "conversation_manager:append_to_conversation")
	defer span.Finish()

	if s.failingMethod == "appendToConversationWith" {
		return false
	}

	// userName is the "other" user (not the one storing)
	val, ok := s.rw.Load(userName)
	if !ok {
		// first message in conversation
		s.rw.Store(userName, []domain.Message{{Author: msgAuthor, Content: msgContent}})
		return true
	}

	conversation, ok := val.([]domain.Message)
	if !ok {
		span.LogFields(log.Error(errors.New("not a conversation stored at Key")))
		return false
	}

	s.rw.Store(userName, append(conversation, domain.Message{Author: msgAuthor, Content: msgContent}))
	return true

}
