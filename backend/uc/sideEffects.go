package uc

import (
	"context"
	"gop2p/domain"
)

// NB : side effects return bool instead of error because we don't want their lower level
// details to leak into usecases, error handling is done at implementation level

// UserStore allows to manage the userStore
type UserStore interface {
	InsertUser(ctx context.Context, login, password string) bool
	GetUserByLoginPassword(ctx context.Context, login, password string) (*domain.User, bool)
	GetUserByLogin(ctx context.Context, login string) (*domain.User, bool)
}

// SessionManager is a struct holding the function types allowing to manage the sessions
type SessionManager interface {
	InsertSession(ctx context.Context, login, address string) bool
	GetSession(ctx context.Context, login string) (*domain.Session, bool)
}

// ConversationManager is used by client to store their conversations with other users
type ConversationManager interface {
	GetConversationWith(ctx context.Context, authorName string) ([]domain.Message, bool)
	AppendToConversationWith(ctx context.Context, userName, msgAuthor, msgContent string) bool
}

// ServerGateway provides client -> server communication
type ServerGateway interface {
	AskSessionToServer(ctx context.Context, from string, to string) (*domain.Session, bool)
}

// ClientGateway provides client -> client communication
type ClientGateway interface {
	SendMsg(ctx context.Context, addr string, msg domain.Message, from string) bool
}
