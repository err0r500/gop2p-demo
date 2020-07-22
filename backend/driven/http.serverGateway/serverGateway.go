package servergateway

import (
	"context"
	"encoding/json"
	"github.com/opentracing/opentracing-go"
	"gop2p/domain"
	"gop2p/driving/api.mux"
	"gop2p/uc"
	"io/ioutil"
	"net/http"
)

type caller struct {
	serverAddress string
	client        *http.Client
	l             uc.Logger
}

func New(serverAddress string, logger uc.Logger) uc.ServerGateway {
	return caller{serverAddress: serverAddress, client: http.DefaultClient, l: logger}
}

func (c caller) AskSessionToServer(ctx context.Context, from string, to string) (*domain.Session, *domain.ErrTechnical) {
	req, err := http.NewRequest(http.MethodGet, "http://"+c.serverAddress+"/sessions/"+to, nil)
	if err != nil {
		c.l.Log(err)
		return nil, &domain.ErrTechnical{}
	}
	req.Header.Set("user", from)

	if span := opentracing.SpanFromContext(ctx); span != nil {
		mux.InjectSpanInReq(span, req)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		c.l.Log(err)
		return nil, &domain.ErrTechnical{}
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.l.Log(err)
		return nil, &domain.ErrTechnical{}
	}

	session := &domain.Session{}
	if err := json.Unmarshal(body, session); err != nil {
		c.l.Log(err)
		return nil, &domain.ErrTechnical{}
	}

	return session, nil
}
