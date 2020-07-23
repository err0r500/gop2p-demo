package mux_test

import (
	"net/http"
	"net/http/httptest"
	"strconv"

	. "github.com/smartystreets/goconvey/convey"
	mux "gop2p/driving/api.mux"
)

type spy struct {
	called int
}

func withServer(router mux.ServerRouter, f func(*httptest.Server)) func() {
	return func() {
		r := http.NewServeMux()
		router.SetRoutes(r)
		s := httptest.NewServer(r)
		defer s.Close()
		f(s)
	}
}

func itRespondsAnEmptyBody(resp *http.Response) {
	Convey("it responds an empty body", func() {
		So(resp.ContentLength, ShouldEqual, 0)
	})
}

func itRespondsWithStatus(status int, resp *http.Response) {
	Convey("it responds "+strconv.Itoa(status), func() {
		So(resp.StatusCode, ShouldEqual, status)
	})
}
