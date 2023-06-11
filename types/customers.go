package types

import "time"

type Customer struct {
	ID       int64     `json:"id"`
	Name     string    `json:"name"`
	Phone    string    `json:"phone"`
	Password *string   `json:"password"`
	Active   bool      `json:"active"`
	Created  time.Time `json:"created"`
	Balance  int64     `json:"balance"`
}
