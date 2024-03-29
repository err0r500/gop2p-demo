package uc_test

import (
	"context"
	"gop2p/domain"
	"gop2p/uc"

	. "github.com/smartystreets/goconvey/convey"
)

func noSessionIsCreated(sm uc.SessionManager, uName string) {
	Convey("no new session is created", func() {
		session, ok := sm.GetSession(context.Background(), uName)
		So(ok, ShouldBeTrue)
		So(session, ShouldBeNil)
	})
}

func aNewSessionIsCreated(sm uc.SessionManager, uName string) {
	Convey("a new session is created", func() {
		session, ok := sm.GetSession(context.Background(), uName)
		So(ok, ShouldBeTrue)
		So(session, ShouldNotBeNil)
	})
}

// errors check
func noErrorReturned(err error) {
	Convey("no error is returned", func() {
		So(err, ShouldBeNil)
	})
}

func errorReturned(err error) {
	Convey("an error is returned", func() {
		So(err, ShouldNotBeNil)
	})
}

func resourceNotFoundErrIsReturned(err error) {
	Convey("a resourceNotFound error is returned", func() {
		So(err, ShouldHaveSameTypeAs, domain.ErrResourceNotFound{})
	})
}

func techErrIsReturned(err error) {
	Convey("a technical error is returned", func() {
		So(err, ShouldHaveSameTypeAs, domain.ErrTechnical{})
	})
}
