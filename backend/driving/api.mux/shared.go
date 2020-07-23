package mux

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"gop2p/domain"
	"log"
	"net/http"
	"strings"
)

func mapDomainErrToHttpCode(ctx context.Context, err error, w http.ResponseWriter) {
	span := opentracing.SpanFromContext(ctx)

	switch err.(type) {
	case nil:
		writeSpanAndHeader(span, w, http.StatusOK)
		return
	case domain.ErrResourceNotFound:
		writeSpanAndHeader(span, w, http.StatusUnauthorized)
		return
	case domain.ErrTechnical:
		writeSpanAndHeader(span, w, http.StatusInternalServerError)
		return
	case domain.ErrMalformed:
		writeSpanAndHeader(span, w, http.StatusBadRequest)
		return
	default:
		writeSpanAndHeader(span, w, http.StatusInternalServerError)
		return
	}
}

const spanHttpStatusKey = "http_status"

func spanHttpOK(span opentracing.Span) {
	span.SetTag(spanHttpStatusKey, http.StatusOK)
}

func writeSpanAndHeader(span opentracing.Span, w http.ResponseWriter, status int) {
	if span != nil {
		span.SetTag(spanHttpStatusKey, status)
	}
	w.WriteHeader(status)
}

func paramAtIndex(r *http.Request, index int) string {
	p := strings.Split(r.URL.Path, "/")
	if len(p) < index {
		return ""
	}
	return p[index]
}

func InjectSpanInReq(span opentracing.Span, req *http.Request) {
	err := opentracing.GlobalTracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header))
	if err != nil {
		log.Println(err)
	}
}

func spanFromReq(spanName string, r *http.Request) opentracing.Span {
	wireContext, err := opentracing.GlobalTracer().Extract(
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(r.Header))
	if err != nil {
		log.Println(spanName, err)
	}

	return opentracing.StartSpan(
		spanName,
		ext.RPCServerOption(wireContext))
}
