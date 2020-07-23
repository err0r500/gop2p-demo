package mux

import (
	"fmt"
	"gop2p/uc"
	"log"
	"net/http"
	"net/url"
)

// ApplicationJSON is the expected content-type, a constant is used to avoid typos
const ApplicationJSON = "application/json"

// ServerRouter is used to map logic (usecases) and http routes
type ServerRouter struct {
	Logic uc.ServerLogic
}

// ClientP2pRouter is the router used by clients for p2p internal communications
type ClientP2pRouter struct {
	Logic uc.ClientP2PLogic
}

// ClientFrontRouter is the router used by clients to allow interactions with the frontend
type ClientFrontRouter struct {
	Logic         uc.ClientFrontLogic
	ServerAddress *url.URL
}

// NewServerRouter initializes the server router
func NewServerRouter(l uc.ServerLogic, port int) {
	mux := http.NewServeMux()
	ServerRouter{
		Logic: l,
	}.SetRoutes(mux)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	fmt.Println("listening on", port)
	log.Fatal(server.ListenAndServe())
}

// NewClientFrontRouter initializes the client frontend router
func NewClientFrontRouter(l uc.ClientFrontLogic, port int, serverAddress string) {
	mux := http.NewServeMux()
	ClientFrontRouter{
		Logic:         l,
		ServerAddress: &url.URL{Host: serverAddress},
	}.SetRoutes(mux)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	fmt.Println("listening on", port)
	log.Fatal(server.ListenAndServe())

}

// NewClientP2pRouter initializes the client p2p router
func NewClientP2pRouter(l uc.ClientP2PLogic, port int) {
	mux := http.NewServeMux()
	ClientP2pRouter{
		Logic: l,
	}.SetRoutes(mux)
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	fmt.Println("listening on", port)
	log.Fatal(server.ListenAndServe())

}

// SetRoutes plugs routes with logic
func (r ServerRouter) SetRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {})
	mux.HandleFunc("/sessions/", serverSessionsHandler(r.Logic))
}

// SetRoutes plugs routes with logic
func (r ClientP2pRouter) SetRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {})
	mux.HandleFunc("/messages/", clientp2pHandler(r.Logic))
}

// SetRoutes plugs routes with logic
func (r ClientFrontRouter) SetRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {})
	mux.HandleFunc("/sessions/", clientFrontSessionsHandler(r.Logic, r.ServerAddress))
	mux.HandleFunc("/conversations/", clientFrontConversationsHandler(r.Logic))
	mux.HandleFunc("/messages/", clientFrontMessagessHandler(r.Logic))
}
