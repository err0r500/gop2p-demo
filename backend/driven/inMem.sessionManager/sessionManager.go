package sessionmanager

import (
	"errors"
	"gop2p/domain"
	"gop2p/uc"
	"sync"
)

type store struct {
	rw *sync.Map
}

// New is the constructor of this in memory implementation of the uc.SessionManager
func New() uc.SessionManager {
	s := store{rw: &sync.Map{}}
	return uc.SessionManager{
		InsertSession: s.InsertSession,
		GetSession:    s.GetSession,
	}
}

func (s store) InsertSession(login, address string) error {
	s.rw.Store(login, domain.Session{Online: true, Address: address})
	return nil
}

func (s store) GetSession(login string) (*domain.Session, error) {
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
