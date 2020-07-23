package mux

import (
	"context"
	"encoding/json"
	"github.com/go-playground/validator"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"gop2p/domain"
	"gop2p/uc"
	"io"
	"net/http"
)

func serverSessionsHandler(serverLogic uc.ServerLogic) func(w http.ResponseWriter, r *http.Request) {
	postHandler := handleStartSession(serverLogic)
	getHandler := handleGetSession(serverLogic)

	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getHandler(w, r)

		case http.MethodPost:
			postHandler(w, r)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

// CreateNewSessionBody is the body of the expected startSession request
type CreateNewSessionBody struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
	Address  string `json:"address"`
}

// FromJSON is the standard json.Unmarshal method
func (nS *CreateNewSessionBody) FromJSON(r io.Reader) error {
	return json.NewDecoder(r).Decode(nS)
}

// Validate is used to check request validity
func (nS *CreateNewSessionBody) Validate() error {
	return validator.New().Struct(nS)
}

func handleStartSession(logic uc.ServerLogic) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		span := spanFromReq("http:handle_start_session", r)
		defer span.Finish()
		ctx := opentracing.ContextWithSpan(context.Background(), span)

		b := CreateNewSessionBody{}
		if err := b.FromJSON(r.Body); err != nil {
			span.LogFields(log.Error(err))
			mapDomainErrToHttpCode(ctx, domain.ErrMalformed{}, w)
			return
		}

		if err := b.Validate(); err != nil {
			span.LogFields(log.Error(err))
			mapDomainErrToHttpCode(ctx, domain.ErrMalformed{}, w)
			return
		}

		address := ""
		if b.Address != "" {
			address = b.Address
		} else {
			address = r.RemoteAddr
		}

		if err := logic.StartSession(ctx, b.Login, b.Password, address); err != nil {
			span.LogFields(log.Error(err))
			mapDomainErrToHttpCode(ctx, err, w)
			return
		}

		spanHttpOK(span)
	}
}

func handleGetSession(logic uc.ServerLogic) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		span := spanFromReq("http:handle_get_session", r)
		defer span.Finish()
		ctx := opentracing.ContextWithSpan(context.Background(), span)

		from := r.Header.Get("user")
		if from == "" {
			mapDomainErrToHttpCode(ctx, domain.ErrUnauthorized{}, w)
			return
		}

		to := paramAtIndex(r, 2) // /sessions/:to
		if to == "" {
			mapDomainErrToHttpCode(ctx, domain.ErrMalformed{}, w)
			return
		}

		s, err := logic.ProvideUserSession(ctx, from, to)
		if err != nil {
			span.LogFields(log.Error(err))
			mapDomainErrToHttpCode(ctx, err, w)
			return
		}

		body, err := json.Marshal(s)
		if err != nil {
			span.LogFields(log.Error(err))
			mapDomainErrToHttpCode(ctx, domain.ErrTechnical{}, w)
			return
		}

		w.Write(body)
		spanHttpOK(span)
	}
}
