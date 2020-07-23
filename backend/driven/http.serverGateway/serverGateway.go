package servergateway

import (
	"context"
	"encoding/json"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"gop2p/domain"
	"gop2p/driving/api.mux"
	"gop2p/uc"
	"io/ioutil"
	"net/http"
)

type caller struct {
	serverAddress string
	client        *http.Client
}

func New(serverAddress string) uc.ServerGateway {
	return caller{serverAddress: serverAddress, client: http.DefaultClient}
}

func (c caller) AskSessionToServer(ctx context.Context, from string, to string) (*domain.Session, *domain.ErrTechnical) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ask_session_to_server")
	defer span.Finish()

	req, err := http.NewRequest(http.MethodGet, "http://"+c.serverAddress+"/sessions/"+to, nil)
	if err != nil {
		span.LogFields(log.Error(err))
		return nil, &domain.ErrTechnical{}
	}
	req.Header.Set("user", from)

	mux.InjectSpanInReq(span, req)

	resp, err := c.client.Do(req)
	if err != nil {
		span.LogFields(log.Error(err))
		return nil, &domain.ErrTechnical{}
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		span.LogFields(log.Error(err))
		return nil, &domain.ErrTechnical{}
	}

	session := &domain.Session{}
	if err := json.Unmarshal(body, session); err != nil {
		span.LogFields(log.Error(err))
		return nil, &domain.ErrTechnical{}
	}

	return session, nil
}
