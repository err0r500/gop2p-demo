package uc

import (
	"gop2p/domain"
	"strconv"
	"strings"
)

// StartSessionInit registers the address where the client can be reached
// returns nil if everything is OK
// usage : first, create an instance by applying the 2 needed functions, then use the returned function wherever you want
func StartSessionInit(
	getUserLoginPass GetUserByLoginPassword,
	insertSess InsertSession,
) StartSession {
	return func(login, password, clientAddress string) error {
		if !validAddress(clientAddress) {
			return domain.ErrMalformed{Details: []string{"the address provided is invalid"}}
		}

		user, err := getUserLoginPass(login, password)
		if err != nil {
			return domain.ErrTechnical{}
		}
		if user == nil {
			return domain.ErrResourceNotFound{}
		}

		if err := insertSess(login, clientAddress); err != nil {
			return domain.ErrTechnical{}
		}
		return nil
	}
}

func validAddress(address string) bool {
	ss := strings.Split(address, ":")
	if len(ss) != 2 {
		return false
	}

	port, err := strconv.Atoi(ss[1])
	if err != nil {
		return false
	}

	return port > 0
}
