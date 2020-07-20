package uc

// StartSession is the signature of the func used by clients to create a new session
type StartSession func(login, password, address string) error

// ServerLogic handles the logic of the central server
type ServerLogic struct {
	StartSession
}

// ClientP2PLogic handles the logic of the central server
type ClientP2PLogic struct{}

// ClientFrontLogic handles the logic exposed to the frontend
type ClientFrontLogic struct{}
