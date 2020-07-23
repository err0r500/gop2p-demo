package mux_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"gop2p/domain"
	"gop2p/uc"

	. "github.com/smartystreets/goconvey/convey"
	mux "gop2p/driving/api.mux"
)

const sessionsPath = "/sessions/"

func TestSessionsPost(t *testing.T) {
	login := "matth"
	password := "dummyPassword"
	address := "address:12345"
	reqBody := newPostSessionReqBody(login, password, address)

	Convey("when /sessions is called with a POST", t, func() {
		spy := new(spy)
		Convey("the usecase is called with the correct params",
			withServer(
				newStartSessionRouterWithParamExpectations(t, spy, login, password, address), func(s *httptest.Server) {
					doPostSessionRequest(s, reqBody)
				}),
		)
		So(spy.called, ShouldEqual, 1)
	})

	Convey("when usecase return is", t, func() {
		Convey("everything workded fine", func() {
			router := setStartSessionUsecaseReturn(nil)
			Convey("then", withServer(router, func(s *httptest.Server) {
				r := doPostSessionRequest(s, reqBody)
				itRespondsWithStatus(http.StatusOK, r)
				itRespondsAnEmptyBody(r)
			}))
		})

		Convey("when the usecase returns a badRequest error", func() {
			router := setStartSessionUsecaseReturn(domain.ErrMalformed{})
			Convey("then", withServer(router, func(s *httptest.Server) {
				r := doPostSessionRequest(s, reqBody)
				itRespondsWithStatus(http.StatusBadRequest, r)
				itRespondsAnEmptyBody(r)
			}))
		})

		Convey("when the usecase returns a userNotFound error", func() {
			router := setStartSessionUsecaseReturn(domain.ErrResourceNotFound{})
			Convey("then", withServer(router, func(s *httptest.Server) {
				r := doPostSessionRequest(s, reqBody)
				itRespondsWithStatus(http.StatusUnauthorized, r)
				itRespondsAnEmptyBody(r)
			}))
		})

		Convey("when the usecase returns a technical error", func() {
			router := setStartSessionUsecaseReturn(domain.ErrTechnical{})
			Convey("then", withServer(router, func(s *httptest.Server) {
				r := doPostSessionRequest(s, reqBody)
				itRespondsWithStatus(http.StatusInternalServerError, r)
				itRespondsAnEmptyBody(r)
			}))
		})

		Convey("special case : if address is missing in request body", func() {
			reqBodyWithoutRMA, err := json.Marshal(mux.CreateNewSessionBody{
				Login:    login,
				Password: password,
			})
			So(err, ShouldBeNil)

			injectedRemoteAddress := "127.0.0.1:54321"
			Convey("req.RemoteAddr header is used as a fallback when calling the usecase",
				withServer(
					newStartSessionRouterWithParamExpectations(t, nil, login, password, injectedRemoteAddress), func(s *httptest.Server) {
						req, err := http.NewRequest(http.MethodPost, s.URL, bytes.NewBuffer(reqBodyWithoutRMA))
						So(err, ShouldBeNil)
						req.RemoteAddr = injectedRemoteAddress

						s.Client().Do(req)
					}),
			)
		})

	})
}

func TestSessionsGet(t *testing.T) {
	from := "alice"
	to := "bob"

	Convey("when /sessions is called with a GET", t, func() {
		spy := new(spy)
		Convey("the usecase is spy with the correct params",
			withServer(
				newProvideSessionRouterWithParamExpectations(t, spy, from, to), func(s *httptest.Server) {
					doGetSessionRequest(s, from, to)
				}),
		)
		So(spy.called, ShouldEqual, 1)
	})
}

func TestSessionsOtherMethods(t *testing.T) {
	Convey("when /session is called with another method", t,
		withServer(mux.ServerRouter{}, func(s *httptest.Server) {
			req, err := http.NewRequest(http.MethodPut, s.URL+sessionsPath, nil)
			So(err, ShouldBeNil)

			r, err := s.Client().Do(req)
			So(err, ShouldBeNil)
			itRespondsWithStatus(http.StatusMethodNotAllowed, r)
			itRespondsAnEmptyBody(r)
		}),
	)

}

func newProvideSessionRouterWithParamExpectations(t *testing.T, spy *spy, from, to string) mux.ServerRouter {
	return mux.ServerRouter{
		Logic: uc.ServerLogic{
			ProvideUserSession: func(_ context.Context, src, dst string) (*domain.Session, error) {
				Convey("provideSession usecase is spy with the right params", t, func() {
					spy.called++
					So(src, ShouldEqual, from)
					So(dst, ShouldEqual, to)
				})
				return nil, nil
			},
		}}
}

func setStartSessionUsecaseReturn(err error) mux.ServerRouter {
	return mux.ServerRouter{Logic: uc.ServerLogic{
		StartSession: func(_ context.Context, _, _, _ string) error {
			return err
		},
	}}
}

// will return a router with the usecase applied.
// it will test if the use case is called with the given params
func newStartSessionRouterWithParamExpectations(t *testing.T, spy *spy, login, password, address string) mux.ServerRouter {
	return mux.ServerRouter{
		Logic: uc.ServerLogic{
			StartSession: func(_ context.Context, l, p, rma string) error {
				Convey("the startSession usecase is called with the right params", t, func() {
					spy.called++
					So(l, ShouldEqual, login)
					So(p, ShouldEqual, password)
					So(rma, ShouldEqual, address)
				})
				return nil
			},
		}}
}

func newPostSessionReqBody(login, password, address string) []byte {
	reqBody, err := json.Marshal(mux.CreateNewSessionBody{
		Login:    login,
		Password: password,
		Address:  address,
	})
	if err != nil {
		log.Fatal(err)
	}
	return reqBody
}

func doPostSessionRequest(s *httptest.Server, reqBody []byte) *http.Response {
	r, err := s.Client().Post(s.URL+sessionsPath, mux.ApplicationJSON, bytes.NewBuffer(reqBody))
	So(err, ShouldBeNil)
	return r
}

func doGetSessionRequest(s *httptest.Server, from, to string) *http.Response {
	req, err := http.NewRequest(http.MethodGet, s.URL+sessionsPath+to, nil)
	So(err, ShouldBeNil)

	req.Header.Add("user", from)

	r, err := s.Client().Do(req)
	So(err, ShouldBeNil)
	return r
}
