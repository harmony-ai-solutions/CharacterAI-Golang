package cai

import "errors"

var (
	// ErrAuthenticationFailed indicates authentication failure
	ErrAuthenticationFailed = errors.New("authentication failed")
	// ErrInvalidResponse indicates an invalid response from the server
	ErrInvalidResponse = errors.New("invalid response from server")
	// Define other custom errors as needed
)
