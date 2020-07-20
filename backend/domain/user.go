package domain

// User is our abstract user (no ORM nor JSON pollution here)
type User struct {
	Login    string
	Password string
}
