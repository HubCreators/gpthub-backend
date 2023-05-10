package domain

import "errors"

var (
	ErrUserAlreadyExists = errors.New("user with such email already exists")
)
