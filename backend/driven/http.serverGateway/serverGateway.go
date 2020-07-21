package servergateway

import (
	"context"
	"encoding/json"
	"github.com/opentracing/opentracing-go"
	"gop2p/domain"
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
	request, err := http.NewRequest(http.MethodGet, "http://"+c.serverAddress+"/sessions/"+to, nil)
	if err != nil {
		c.l.Log(err)
		return nil, &domain.ErrTechnical{}
	}
	request.Header.Set("user", from)

	if span := opentracing.SpanFromContext(ctx); span != nil {
		if err := opentracing.GlobalTracer().Inject(
			span.Context(),
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(request.Header)); err != nil {
			c.l.Log(err)
		}
	}

	resp, err := c.client.Do(request)
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
