package types

import "errors"

var ErrNoSuchUser = errors.New("no such user")
var ErrInvalidPassword = errors.New("invalid password")
var ErrNoSuch = errors.New("more than you have")
