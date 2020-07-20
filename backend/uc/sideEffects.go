package uc

import "gop2p/domain"

// InsertUser is the signature of the func to insert a user in the userStore
type InsertUser func(login, password string) error

// GetUserByLoginPassword is the signature of the func to get a user with matching login and password
type GetUserByLoginPassword func(login, password string) (*domain.User, error)

// UserStore is a struct holding the function types allowing to manage the userStore
// NB, it's a struct because only the function signatures are provided, not their names (it's not methods)
// you can think of this struct as an interface
type UserStore struct {
	InsertUser
	GetUserByLoginPassword
}

// InsertSession is the signature of the func to insert a new session in the session manager
type InsertSession func(login, address string) error

// GetSession is the signature of the func to get the session details for a given login
type GetSession func(login string) (*domain.Session, error)

// SessionManager is a struct holding the function types allowing to manage the sessions
type SessionManager struct {
	InsertSession
	GetSession
}
