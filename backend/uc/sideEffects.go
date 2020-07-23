package uc

import (
	"context"
	"gop2p/domain"
)

// UserStore allows to manage the userStore
type UserStore interface {
	InsertUser(ctx context.Context, login, password string) error
	GetUserByLoginPassword(ctx context.Context, login, password string) (*domain.User, error)
	GetUserByLogin(ctx context.Context, login string) (*domain.User, error)
}

// SessionManager is a struct holding the function types allowing to manage the sessions
type SessionManager interface {
	InsertSession(ctx context.Context, login, address string) error
	GetSession(ctx context.Context, login string) (*domain.Session, error)
}

// ConversationManager is used by client to store their conversations with other users
type ConversationManager interface {
	GetConversationWith(ctx context.Context, authorName string) ([]domain.Message, error)
	AppendToConversationWith(ctx context.Context, userName, msgAuthor, msgContent string) error
}

// ServerGateway provides client -> server communication
type ServerGateway interface {
	AskSessionToServer(ctx context.Context, from string, to string) (*domain.Session, error)
}

// ClientGateway provides client -> client communication
type ClientGateway interface {
	SendMsg(ctx context.Context, addr string, msg domain.Message, from string) error
}
