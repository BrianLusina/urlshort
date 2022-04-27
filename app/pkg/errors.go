package pkg

import "errors"

var (
	ErrServerError = errors.New("unexpected error encountered in server side")
	ErrInvalidUrl  = errors.New("invalid url")
)
