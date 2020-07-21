package uc

import (
	"github.com/opentracing/opentracing-go"
	"gop2p/domain"
	"strconv"
	"strings"
)

// ServerLogic handles the logic of the central server, we use a struct in order to be able to easily change
// implementations in tests and because having several implementation is not very likely
type ServerLogic struct {
	StartSession       func(login, password, address string) error
	ProvideUserSession func(srcLogin, dstLogin string) (*domain.Session, error)
}

type serverInteractor struct {
	uS UserStore
	sM SessionManager
}

func NewServerLogic(uS UserStore, sM SessionManager) ServerLogic {
	i := serverInteractor{
		uS,
		sM,
	}
	return ServerLogic{
		StartSession:       i.StartSession,
		ProvideUserSession: i.ProvideUserSession,
	}
}

// StartSessionInit registers the address where the client can be reached
// returns nil if everything is OK
func (i serverInteractor) StartSession(login, password, clientAddress string) error {
	span := opentracing.GlobalTracer().StartSpan("uc:start_new_session")
	defer span.Finish()

	if !validAddress(clientAddress) {
		return domain.ErrMalformed{Details: []string{"the address provided is invalid"}}
	}

	user, err := i.uS.GetUserByLoginPassword(login, password)
	if err != nil {
		return domain.ErrTechnical{}
	}
	if user == nil {
		return domain.ErrResourceNotFound{}
	}

	if err := i.sM.InsertSession(login, clientAddress); err != nil {
		return domain.ErrTechnical{}
	}
	return nil
}

// ProvideUserSessionInit allows a client to get the session details of another one
func (i serverInteractor) ProvideUserSession(srcLogin, dstLogin string) (*domain.Session, error) {
	span := opentracing.GlobalTracer().StartSpan("uc:provide_user_session")
	defer span.Finish()

	// only users with a session can ask for another user session
	u, err := i.uS.GetUserByLogin(srcLogin)
	if err != nil {
		return nil, domain.ErrTechnical{}
	}
	if u == nil {
		return nil, domain.ErrUnauthorized{}
	}

	s, err := i.sM.GetSession(dstLogin)
	if err != nil {
		return nil, domain.ErrTechnical{}
	}
	if s == nil {
		return nil, domain.ErrResourceNotFound{}
	}

	return s, nil
}

func validAddress(address string) bool {
	ss := strings.Split(address, ":")
	if len(ss) != 2 {
		return false
	}

	port, err := strconv.Atoi(ss[1])
	if err != nil {
		return false
	}

	return port > 0
}