package mux_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	mux "gop2p/driving/api.mux"
)

func TestHealth(t *testing.T) {
	Convey("when /healthz is called", t,
		withServer(mux.ServerRouter{}, func(s *httptest.Server) {
			resp, err := s.Client().Get(s.URL + "/healthz")
			So(err, ShouldBeNil)

			Convey("it responds 200", func() {
				So(resp.StatusCode, ShouldEqual, http.StatusOK)
			})
			Convey("it responds an empty body", func() {
				So(resp.ContentLength, ShouldEqual, 0)
			})
		}),
	)
}
