package userstore

import (
	"errors"
	"gop2p/domain"
	"gop2p/uc"
	"sync"
)

type store struct {
	rw *sync.Map
}

// New is the constructor of this in memory implementation of the uc.UserStore
func New() uc.UserStore {
	s := &store{rw: &sync.Map{}}
	return uc.UserStore{
		InsertUser:             s.InsertUser,
		GetUserByLoginPassword: s.GetUserByLoginPassword,
	}
}

func (s store) InsertUser(login, password string) error {
	// obviously we wouldn't store users with plain-text password in real-life
	s.rw.Store(login, domain.User{Login: login, Password: password})
	return nil
}

func (s store) GetUserByLoginPassword(login, password string) (*domain.User, error) {
	val, ok := s.rw.Load(login)
	if !ok {
		return nil, nil
	}

	user, ok := val.(domain.User)
	if !ok {
		return nil, errors.New("not a user stored at key")
	}
	if user.Password != password {
		return nil, nil
	}

	return &user, nil
}
