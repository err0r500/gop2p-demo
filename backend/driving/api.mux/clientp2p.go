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

func clientp2pHandler(logic uc.ClientP2PLogic) func(w http.ResponseWriter, r *http.Request) {
	handler := handleMessageReceived(logic)

	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handler(w, r)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

// PostMessageBody is the body of the expected handleNewMessage request
type PostMessageBody struct {
	Message string `json:"message" validate:"required"`
}

// FromJSON is the standard json.Unmarshal method
func (nS *PostMessageBody) FromJSON(r io.Reader) error {
	return json.NewDecoder(r).Decode(nS)
}

// Validate is used to check request validity
func (nS *PostMessageBody) Validate() error {
	return validator.New().Struct(nS)
}

func handleMessageReceived(logic uc.ClientP2PLogic) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		span := spanFromReq("http:p2p_message_received", r)
		defer span.Finish()
		ctx := opentracing.ContextWithSpan(context.Background(), span)

		from := r.Header.Get("user")
		if from == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		b := PostMessageBody{}
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

		if err := logic.HandleMessageReceived(ctx, b.Message, domain.User{Login: from}); err != nil {
			span.LogFields(log.Error(err))
			mapDomainErrToHttpCode(ctx, err, w)
		}

		spanHttpOK(span)
	}
}
