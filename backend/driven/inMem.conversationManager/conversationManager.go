package conversationmanager

import (
	"errors"
	"gop2p/domain"
	"gop2p/uc"
	"log"
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

func (s store) GetConversationWith(authorName string) ([]domain.Message, error) {
	if s.failingMethod == "getConversationWith" {
		return nil, errors.New("woops")
	}

	val, ok := s.rw.Load(authorName)
	if !ok {
		return nil, nil
	}

	conversation, ok := val.([]domain.Message)
	if !ok {
		return nil, errors.New("not a conversation stored at key")
	}
	return conversation, nil
}

func (s store) AppendToConversationWith(userName, msgAuthor, msgContent string) error {
	if s.failingMethod == "appendToConversationWith" {
		return errors.New("woops")
	}

	// userName is the "other" user (not the one storing)
	val, ok := s.rw.Load(userName)
	if !ok {
		log.Println("first message in conversation")
		// first message in conversation
		s.rw.Store(userName, []domain.Message{{Author: msgAuthor, Content: msgContent}})
		return nil
	}

	conversation, ok := val.([]domain.Message)
	if !ok {
		err := errors.New("not a conversation stored at key")
		log.Println(err)
		return err
	}

	s.rw.Store(userName, append(conversation, domain.Message{Author: msgAuthor, Content: msgContent}))
	return nil
}
