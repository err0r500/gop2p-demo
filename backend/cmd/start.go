package cmd

import (
	"fmt"

	"gop2p/uc"

	mux "gop2p/driving/api.mux"
)

func startInClientMode(apiPort, p2pPort int) {
	fmt.Println("== RUNNING IN CLIENT MODE ==")

	// in client mode we have 2 servers running :
	go func() {
		uc := uc.ClientFrontLogic{}
		// handles client's frontend traffic
		mux.NewClientFrontRouter(uc, apiPort)
	}()

	// handles p2p traffic
	uc := uc.ClientP2PLogic{}
	mux.NewClientP2pRouter(uc, p2pPort)
}

func startInServerMode(apiPort int) {
	fmt.Println("== RUNNING IN SERVER MODE ==")
	uc := uc.ServerLogic{}
	mux.NewServerRouter(uc, apiPort)
}
