package errors

import err "errors"

var (
	ErrServerError = err.New("unexpected error encountered in server side")
	ErrInvalidUrl  = err.New("invalid url")
	ErrUrlNotFound = err.New("url not found")
	ErrUrlExists   = err.New("url already exists")
)
