package sessionmanager

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
func New() uc.SessionManager {
	return store{rw: &sync.Map{}}
}

type FailingSessionManager interface {
	uc.SessionManager
	InjectErrorAt(failingMethod string)
}

func NewFailable() FailingSessionManager {
	return &store{rw: &sync.Map{}, failingMethod: ""}
}

func (s *store) InjectErrorAt(failingMethod string) {
	s.failingMethod = failingMethod
}

func (s store) InsertSession(login, address string) error {
	if s.failingMethod == "insertSession" {
		return errors.New("woops")
	}
	log.Println("storing session", login, address)

	s.rw.Store(login, domain.Session{Online: true, Address: address})
	return nil
}

func (s store) GetSession(login string) (*domain.Session, error) {
	if s.failingMethod == "getSession" {
		return nil, errors.New("woops")
	}

	val, ok := s.rw.Load(login)
	if !ok {
		return nil, nil
	}

	session, ok := val.(domain.Session)
	if !ok {
		return nil, errors.New("not a session stored at key")
	}

	return &session, nil
}
