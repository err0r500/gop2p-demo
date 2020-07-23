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
	"net/http/httputil"
	"net/url"
)

func clientFrontSessionsHandler(logic uc.ClientFrontLogic, serverAddress *url.URL) func(w http.ResponseWriter, r *http.Request) {
	handleSuccessfulRegistration := func(ctx context.Context, username string) func(resp *http.Response) (err error) {
		return func(resp *http.Response) (err error) {
			span, ctx := opentracing.StartSpanFromContext(ctx, "post_session_response")
			defer span.Finish()

			if resp.StatusCode == http.StatusOK {
				if err := logic.NewSessionRegistered(ctx, username); err != nil {
					return err
				}
				spanHttpOK(span)
				return
			}
			span.LogFields(log.Message(resp.Status))
			return
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			span := opentracing.GlobalTracer().StartSpan("http:post_session")
			defer span.Finish()
			ctx := opentracing.ContextWithSpan(context.Background(), span)

			proxy := httputil.ReverseProxy{
				Director: func(req *http.Request) {
					InjectSpanInReq(span, req)
					req.Header.Add("X-Forwarded-Host", r.Host)
					req.Header.Add("X-Origin-Host", r.Host)
					req.URL.Scheme = "http"
					req.URL.Host = serverAddress.Host
					req.Host = serverAddress.Host
				},
			}

			// since we just forward the request, we'd like not to access the body so we use the "user" header to get
			// the caller username, this is clearly not super smart
			proxy.ModifyResponse = handleSuccessfulRegistration(ctx, r.Header.Get("user"))
			proxy.ServeHTTP(w, r)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func clientFrontMessagessHandler(logic uc.ClientFrontLogic) func(w http.ResponseWriter, r *http.Request) {
	handler := handleSendMessageToOtherClient(logic)

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
type SendNewMessageBody struct {
	Message string `json:"message" validate:"required"`
	To      string `json:"to" validate:"required"`
}

// FromJSON is the standard json.Unmarshal method
func (nS *SendNewMessageBody) FromJSON(r io.Reader) error {
	return json.NewDecoder(r).Decode(nS)
}

// Validate is used to check request validity
func (nS *SendNewMessageBody) Validate() error {
	return validator.New().Struct(nS)
}

func handleSendMessageToOtherClient(logic uc.ClientFrontLogic) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		span := opentracing.GlobalTracer().StartSpan("http:post_messages")
		defer span.Finish()
		ctx := opentracing.ContextWithSpan(context.Background(), span)

		b := SendNewMessageBody{}
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

		if err := logic.SendMessageToOtherClient(ctx, b.To, b.Message); err != nil {
			span.LogFields(log.Error(err))
			mapDomainErrToHttpCode(ctx, err, w)
		}
		spanHttpOK(span)
	}
}

func clientFrontConversationsHandler(logic uc.ClientFrontLogic) func(w http.ResponseWriter, r *http.Request) {
	handler := handleGetConversationWith(logic)

	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler(w, r)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func handleGetConversationWith(logic uc.ClientFrontLogic) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		span := opentracing.GlobalTracer().StartSpan("http:get_conversations")
		defer span.Finish()
		ctx := opentracing.ContextWithSpan(context.Background(), span)

		with := paramAtIndex(r, 2) // /conversations/:to
		if with == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		messages, err := logic.GetConversationWith(ctx, with)
		if err != nil {
			span.LogFields(log.Error(err))
			mapDomainErrToHttpCode(ctx, err, w)
		}

		body, err := json.Marshal(messages)
		if err != nil {
			span.LogFields(log.Error(err))
			mapDomainErrToHttpCode(ctx, err, w)
		}

		w.Write(body)
		spanHttpOK(span)
	}
}
