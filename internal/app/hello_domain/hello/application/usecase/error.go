package usecase

import "errors"

var (
	ErrNameRequired    = errors.New("name is required")
	ErrGreetingFailed  = errors.New("failed to create greeting")
)
