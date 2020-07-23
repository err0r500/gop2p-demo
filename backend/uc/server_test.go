package uc_test

import (
	"context"
	"gop2p/domain"
	"gop2p/uc"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	sessionManager "gop2p/driven/inMem.sessionManager"
	userStore "gop2p/driven/inMem.userStore"
)

// cleanServerLogic provides the startSession usecase function with fresh stores
func cleanServerLogic() (uc.UserStore, uc.SessionManager, uc.ServerLogic) {
	us := userStore.New()
	sm := sessionManager.New()
	return us, sm, uc.NewServerLogic(us, sm)
}

func TestStartSession(t *testing.T) {
	uName := "alice"
	uPswd := "alicePass"
	address := "alice-machine:1234"
	ctx := context.Background()

	Convey("given a known user", t, func() {
		uS, sM, sI := cleanServerLogic()
		noErrorReturned(uS.InsertUser(ctx, uName, uPswd))

		Convey("when he attempts to create a new session with valid creds & address", func() {
			ucRet := sI.StartSession(ctx, uName, uPswd, address)
			aNewSessionIsCreated(sM, uName)
			noErrorReturned(ucRet)
		})

		Convey("same happy case but with invalid address", func() {
			Convey("must have 2 part like host:port", func() {
				ucRet := sI.StartSession(ctx, uName, uPswd, "anywhere")
				noSessionIsCreated(sM, uName)
				errorReturned(ucRet)
			})
			Convey("port must be an int", func() {
				ucRet := sI.StartSession(ctx, uName, uPswd, "anywhere:abc")
				noSessionIsCreated(sM, uName)
				errorReturned(ucRet)
			})
			Convey("port must be larger than 0", func() {
				ucRet := sI.StartSession(ctx, uName, uPswd, "anywhere:0")
				noSessionIsCreated(sM, uName)
				errorReturned(ucRet)
			})
		})

		Convey("when another, unknown, user attempts to login", func() {
			unknownUsername := "unknownUsername"
			ucRet := sI.StartSession(ctx, unknownUsername, uPswd, address)
			noSessionIsCreated(sM, unknownUsername)
			resourceNotFoundErrIsReturned(ucRet)
		})

		Convey("when the same user, with the wrong password attempts to login", func() {
			wrongPassword := "wrongPass"
			ucRet := sI.StartSession(ctx, uName, wrongPassword, address)
			noSessionIsCreated(sM, uName)
			resourceNotFoundErrIsReturned(ucRet)
		})
	})

	Convey("when everything should go fine", t, func() {
		us := userStore.NewFailable()
		sm := sessionManager.NewFailable()
		noErrorReturned(us.InsertUser(ctx, uName, uPswd))

		Convey("if a tech error happens with the uS", func() {
			us.InjectErrorAt("getUserByLogicPassword")
			ucRet := uc.NewServerLogic(us, sessionManager.New()).
				StartSession(ctx, uName, uPswd, address)

			noSessionIsCreated(sm, uName)
			techErrIsReturned(ucRet)
		})

		Convey("if a tech error happens with the sessionStore", func() {
			sm.InjectErrorAt("insertSession")

			ucRet := uc.NewServerLogic(us, sm).
				StartSession(ctx, uName, uPswd, address)

			noSessionIsCreated(sm, uName)
			techErrIsReturned(ucRet)
		})
	})
}

func TestGetUserSession(t *testing.T) {
	bobName := "bob"
	bobAddr := "bob:1234"
	aliceName := "alice"
	aliceAddr := "alice:2345"
	ctx := context.Background()

	Convey("given 2 connected users", t, func() {
		userStore, sessionManager, sI := cleanServerLogic()
		So(userStore.InsertUser(ctx, bobName, "pass"), ShouldBeNil)
		So(userStore.InsertUser(ctx, aliceName, "pass"), ShouldBeNil)
		So(sessionManager.InsertSession(ctx, bobName, bobAddr), ShouldBeNil)
		So(sessionManager.InsertSession(ctx, aliceName, aliceAddr), ShouldBeNil)

		Convey("they are able to get each other's session", func() {
			aliceSession, err := sI.ProvideUserSession(ctx, bobName, aliceName)
			So(err, ShouldBeNil)
			So(aliceSession, ShouldNotBeNil)
			So(aliceSession.Address, ShouldEqual, aliceAddr)

			bobSess, err := sI.ProvideUserSession(ctx, aliceName, bobName)
			So(err, ShouldBeNil)
			So(bobSess, ShouldNotBeNil)
			So(bobSess.Address, ShouldEqual, bobAddr)

		})

		Convey("a user has to be connected in order to retrieve a session", func() {
			s, err := sI.ProvideUserSession(ctx, "unknown", aliceName)
			Convey("an Unauthorized error is returned", func() {
				So(err, ShouldHaveSameTypeAs, domain.ErrUnauthorized{})

				So(s, ShouldBeNil)
			})
		})

		Convey("if the queried account has no session", func() {
			s, err := sI.ProvideUserSession(ctx, aliceName, "unknown")
			Convey("a notFoundResource Error is returned", func() {
				So(err, ShouldHaveSameTypeAs, domain.ErrResourceNotFound{})
				So(s, ShouldBeNil)
			})
		})
	})

	Convey("when everything should go fine", t, func() {
		us := userStore.NewFailable()
		sm := sessionManager.NewFailable()
		So(us.InsertUser(ctx, bobName, "pass"), ShouldBeNil)
		So(us.InsertUser(ctx, aliceName, "pass"), ShouldBeNil)
		So(sm.InsertSession(ctx, bobName, bobAddr), ShouldBeNil)
		So(sm.InsertSession(ctx, aliceName, aliceAddr), ShouldBeNil)

		Convey("but a tech error happens when attempting to getUserByLogin", func() {
			us.InjectErrorAt("getUserByLogin")
			s, err := uc.NewServerLogic(us, sm).ProvideUserSession(ctx, aliceName, bobName)
			techErrIsReturned(err)
			So(s, ShouldBeNil)
		})

		Convey("but a tech error happens when attempting to getSession", func() {
			sm.InjectErrorAt("getSession")
			s, err := uc.NewServerLogic(us, sm).ProvideUserSession(ctx, aliceName, bobName)
			techErrIsReturned(err)
			So(s, ShouldBeNil)
		})
	})
}
