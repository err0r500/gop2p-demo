package sessionmanager

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
func New() uc.SessionManager {
	return store{rw: &sync.Map{}}
}

type FailingSessionManager interface {
	uc.SessionManager
	InjectErrorAt(failingMethod string)
}

// NewFailable is just for testing purposes
func NewFailable() FailingSessionManager {
	return &store{rw: &sync.Map{}, failingMethod: ""}
}

func (s *store) InjectErrorAt(failingMethod string) {
	s.failingMethod = failingMethod
}

func (s store) InsertSession(ctx context.Context, login, address string) bool {
	span, ctx := opentracing.StartSpanFromContext(ctx, "session_manager:insert_session")
	defer span.Finish()

	if s.failingMethod == "insertSession" {
		return false
	}

	s.rw.Store(login, domain.Session{Online: true, Address: address})
	return true
}

func (s store) GetSession(ctx context.Context, login string) (*domain.Session, bool) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "session_manager:get_session")
	defer span.Finish()

	if s.failingMethod == "getSession" {
		return nil, false
	}

	val, ok := s.rw.Load(login)
	if !ok {
		return nil, true
	}

	session, ok := val.(domain.Session)
	if !ok {
		err := errors.New("not a session stored at Key")
		span.LogFields(log.Error(err))
		return nil, true
	}

	return &session, true
}
