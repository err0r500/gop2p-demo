package domain

// Session is used to know if a client is online
// and the address he can be reached at
type Session struct {
	Online  bool
	Address string
}
