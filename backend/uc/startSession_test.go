package uc_test

import (
	"gop2p/uc"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	sessionManager "gop2p/driven/inMem.sessionManager"
	userStore "gop2p/driven/inMem.userStore"
)

// clean provides the startSession usecase function with fresh stores
func clean() (uc.UserStore, uc.SessionManager, uc.StartSession) {
	us := userStore.New()
	sm := sessionManager.New()

	return us, sm, uc.StartSessionInit(
		us.GetUserByLoginPassword,
		sm.InsertSession,
	)
}

func TestStartSession(t *testing.T) {
	uName := "matth"
	uPswd := "pass"
	address := "anywhere:1234"

	Convey("given a known user", t, func() {
		userStore, sessionManager, startSession := clean()
		noErrorReturned(userStore.InsertUser(uName, uPswd))

		Convey("when he attempts to create a new session with valid creds & address", func() {
			ucRet := startSession(uName, uPswd, address)
			aNewSessionIsCreated(sessionManager, uName)
			noErrorReturned(ucRet)
		})

		Convey("same happy case but with invalid address", func() {
			Convey("must have 2 part like host:port", func() {
				ucRet := startSession(uName, uPswd, "anywhere")
				noSessionIsCreated(sessionManager, uName)
				errorReturned(ucRet)
			})
			Convey("port must be an int", func() {
				ucRet := startSession(uName, uPswd, "anywhere:abc")
				noSessionIsCreated(sessionManager, uName)
				errorReturned(ucRet)
			})
			Convey("port must be larger than 0", func() {
				ucRet := startSession(uName, uPswd, "anywhere:0")
				noSessionIsCreated(sessionManager, uName)
				errorReturned(ucRet)
			})
		})

		Convey("when another, unknown, user attempts to login", func() {
			unknownUsername := "unknownUsername"
			ucRet := startSession(unknownUsername, uPswd, address)
			noSessionIsCreated(sessionManager, unknownUsername)
			resourceNotFoundErrIsReturned(ucRet)
		})

		Convey("when the same user, with the wrong password attempts to login", func() {
			wrongPassword := "wrongPass"
			ucRet := startSession(uName, wrongPassword, address)
			noSessionIsCreated(sessionManager, uName)
			resourceNotFoundErrIsReturned(ucRet)
		})

		Convey("if a tech error happens with the userStore", func() {
			ucRet := uc.StartSessionInit(
				failingGetUser,
				sessionManager.InsertSession,
			)(uName, uPswd, address)
			noSessionIsCreated(sessionManager, uName)
			techErrIsReturned(ucRet)
		})

		Convey("if a tech error happens with the sessionStore", func() {
			ucRet := uc.StartSessionInit(
				userStore.GetUserByLoginPassword,
				failingInsertSession,
			)(uName, uPswd, address)
			noSessionIsCreated(sessionManager, uName)
			techErrIsReturned(ucRet)
		})
	})
}
