package clientgateway

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"gop2p/domain"
	"gop2p/driving/api.mux"
	"gop2p/uc"
	"net/http"
)

type caller struct {
	client *http.Client
}

func New() uc.ClientGateway {
	return caller{client: http.DefaultClient}
}

func (c caller) SendMsg(ctx context.Context, addr string, msg domain.Message, from string) bool {
	span, ctx := opentracing.StartSpanFromContext(ctx, "http:send_message")
	defer span.Finish()

	reqBody, err := json.Marshal(mux.PostMessageBody{Message: msg.Content})
	if err != nil {
		span.LogFields(log.Error(err))
		return false
	}

	req, err := http.NewRequest(http.MethodPost, "http://"+addr+"/messages/", bytes.NewBuffer(reqBody))
	if err != nil {
		span.LogFields(log.Error(err))
		return false
	}
	req.Header.Set("user", from)

	mux.InjectSpanInReq(span, req)

	resp, err := c.client.Do(req)
	if err != nil {
		span.LogFields(log.Error(err))
		return false
	}
	if resp.StatusCode != http.StatusOK {
		span.LogFields(log.Error(err))
		return false
	}

	return true
}
