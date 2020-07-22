package clientgateway

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/opentracing/opentracing-go"
	"gop2p/domain"
	"gop2p/driving/api.mux"
	"gop2p/uc"
	"log"
	"net/http"
)

type caller struct {
	client *http.Client
}

func New() uc.ClientGateway {
	return caller{client: http.DefaultClient}
}

func (c caller) SendMsg(ctx context.Context, addr string, msg domain.Message, from string) error {
	reqBody, err := json.Marshal(mux.PostMessageBody{Message: msg.Content})
	if err != nil {
		log.Println(err)
		return err
	}

	req, err := http.NewRequest(http.MethodPost, "http://"+addr+"/messages/", bytes.NewBuffer(reqBody))
	if err != nil {
		log.Println(err)
		return err
	}
	req.Header.Set("user", from)

	if span := opentracing.SpanFromContext(ctx); span != nil {
		mux.InjectSpanInReq(span, req)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		log.Println(err)
		return err
	}
	if resp.StatusCode != http.StatusOK {
		log.Println(resp)
		return domain.ErrTechnical{}
	}

	return nil
}
