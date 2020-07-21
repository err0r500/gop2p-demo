package userstore

import (
	"errors"
	"gop2p/domain"
	"gop2p/uc"
	"sync"
)

type store struct {
	rw            *sync.Map
	failingMethod string
}

// New is the constructor of this in memory implementation of the uc.UserStore
func New() uc.UserStore {
	return &store{rw: &sync.Map{}}
}

type FailingStore interface {
	uc.UserStore
	InjectErrorAt(failingMethod string)
}

func NewFailable() FailingStore {
	return &store{rw: &sync.Map{}, failingMethod: ""}
}

func (s *store) InjectErrorAt(failingMethod string) {
	s.failingMethod = failingMethod
}

func (s store) InsertUser(login, password string) error {
	if s.failingMethod == "insertUser" {
		return errors.New("woops")
	}

	// obviously we wouldn't store users with plain-text password in real-life
	s.rw.Store(login, domain.User{Login: login, Password: password})
	return nil
}

func (s store) GetUserByLoginPassword(login, password string) (*domain.User, error) {
	if s.failingMethod == "getUserByLogicPassword" {
		return nil, errors.New("woops")
	}

	val, ok := s.rw.Load(login)
	if !ok {
		return nil, nil
	}

	user, ok := val.(domain.User)
	if !ok {
		return nil, errors.New("not a user stored at Key")
	}
	if user.Password != password {
		return nil, nil
	}

	return &user, nil
}

func (s store) GetUserByLogin(login string) (*domain.User, error) {
	if s.failingMethod == "getUserByLogin" {
		return nil, errors.New("woops")
	}

	val, ok := s.rw.Load(login)
	if !ok {
		return nil, nil
	}

	user, ok := val.(domain.User)
	if !ok {
		return nil, errors.New("not a user stored at Key")
	}

	return &user, nil
}
