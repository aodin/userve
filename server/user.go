package server

import ()

type User interface {
	Name() string
	IsAuthenticated() bool
}

// AnonUser is a anonymous user that can never be considered authenticated
type AnonUser struct{}

func (u *AnonUser) Name() string {
	return ""
}

// IsAuthenticated will always return false for AnonUser instances
func (u *AnonUser) IsAuthenticated() bool {
	return false
}

type AuthUser struct {
	id   int64
	name string
}

func (u *AuthUser) Name() string {
	return u.name
}

// IsAuthenticated will always return false for AnonUser instances
func (u *AuthUser) IsAuthenticated() bool {
	return u.id != 0
}
