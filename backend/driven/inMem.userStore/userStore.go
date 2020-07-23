package userstore

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

func (s store) InsertUser(ctx context.Context, login, password string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "user_store:insert_user")
	defer span.Finish()

	if s.failingMethod == "insertUser" {
		return errors.New("woops")
	}

	// obviously we wouldn't store users with plain-text password in real-life
	s.rw.Store(login, domain.User{Login: login, Password: password})
	return nil
}

func (s store) GetUserByLoginPassword(ctx context.Context, login, password string) (*domain.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "user_store:get_user_by_login_pass")
	defer span.Finish()

	if s.failingMethod == "getUserByLogicPassword" {
		return nil, errors.New("woops")
	}

	val, ok := s.rw.Load(login)
	if !ok {
		return nil, nil
	}

	user, ok := val.(domain.User)
	if !ok {
		err := errors.New("not a user stored at Key")
		span.LogFields(log.Error(err))
		return nil, err
	}
	if user.Password != password {
		span.LogFields(log.Event("passwords don't match"))
		return nil, nil
	}

	return &user, nil
}

func (s store) GetUserByLogin(ctx context.Context, login string) (*domain.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "user_store:get_user_by_login")
	defer span.Finish()

	if s.failingMethod == "getUserByLogin" {
		return nil, errors.New("woops")
	}

	val, ok := s.rw.Load(login)
	if !ok {
		return nil, nil
	}

	user, ok := val.(domain.User)
	if !ok {
		err := errors.New("not a user stored at Key")
		span.LogFields(log.Error(err))
		return nil, err
	}

	return &user, nil
}
