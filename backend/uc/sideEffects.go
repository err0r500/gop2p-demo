package uc

import (
	"context"
	"gop2p/domain"
)

// UserStore allows to manage the userStore
type UserStore interface {
	InsertUser(login, password string) error
	GetUserByLoginPassword(login, password string) (*domain.User, error)
	GetUserByLogin(string) (*domain.User, error)
}

// SessionManager is a struct holding the function types allowing to manage the sessions
type SessionManager interface {
	InsertSession(login, address string) error
	GetSession(login string) (*domain.Session, error)
}

// ConversationManager is used by client to store their conversations with other users
type ConversationManager interface {
	GetConversationWith(authorName string) ([]domain.Message, error)
	AppendToConversationWith(userName, msgAuthor, msgContent string) error
}

// ServerGateway provides client -> server communication
type ServerGateway interface {
	AskSessionToServer(ctx context.Context, from string, to string) (*domain.Session, *domain.ErrTechnical)
}

// ClientGateway provides client -> client communication
type ClientGateway interface {
	SendMsg(addr string, msg domain.Message, from string) error
}

type Logger interface {
	Log(...interface{})
}
