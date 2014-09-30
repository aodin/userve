package server

import ()

type User interface {
	IsAuthenticated() bool
}

// AnonUser is a anonymous user that can never be considered authenticated
type AnonUser struct{}

// IsAuthenticated will always return false for AnonUser instances
func (u *AnonUser) IsAuthenticated() bool {
	return false
}
