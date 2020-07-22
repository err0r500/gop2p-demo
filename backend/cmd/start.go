package cmd

import (
	"fmt"
	"gop2p/driven/http.clientGateway"
	"gop2p/driven/http.serverGateway"
	"gop2p/driven/inMem.conversationManager"
	"gop2p/driven/inMem.sessionManager"
	"gop2p/driven/inMem.userStore"
	"io"

	"gop2p/uc"

	mux "gop2p/driving/api.mux"
	"log"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-lib/metrics"

	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
)

type simpleLogger struct{}

func (simpleLogger) Log(ll ...interface{}) {
	log.Println(ll...)
}

func setTracer() (opentracing.Tracer, io.Closer) {
	cfg, err := jaegercfg.FromEnv()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize tracer with a logger and a metrics factory
	tracer, closer, err := cfg.NewTracer(
		jaegercfg.Logger(jaegerlog.StdLogger),
		jaegercfg.Metrics(metrics.NullFactory),
	)
	if err != nil {
		log.Fatal(err)
	}
	return tracer, closer
}

func startInClientMode(apiPort, p2pPort int, serverAddress string) {
	fmt.Println("== RUNNING IN CLIENT MODE ==")

	tracer, closer := setTracer()
	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()

	// in client mode we have 2 servers running :
	cm := conversationmanager.New()

	go func(cm uc.ConversationManager) {
		// handles client's frontend traffic
		mux.NewClientFrontRouter(
			uc.NewClientFrontLogic(
				cm,
				servergateway.New(serverAddress, simpleLogger{}),
				clientgateway.New(),
			),
			apiPort,
			serverAddress,
		)
	}(cm)

	// handles p2p traffic
	mux.NewClientP2pRouter(uc.NewClientP2pLogic(cm), p2pPort)
}

func startInServerMode(apiPort int) {
	fmt.Println("== RUNNING IN SERVER MODE ==")

	tracer, closer := setTracer()
	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()

	us := userstore.New()

	// we just add 2 users for testing
	us.InsertUser("alice", "pass")
	us.InsertUser("bob", "pass")

	mux.NewServerRouter(
		uc.NewServerLogic(
			us,
			sessionmanager.New(),
		),
		apiPort,
	)
}
