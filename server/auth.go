package server

import (
	"fmt"
	"net/http"

	"github.com/aodin/volta/auth"
	"github.com/aodin/volta/config"
)

// TODO Move to actual auth package
func AuthLogin(w http.ResponseWriter, username, password string, sessions auth.SessionManager, users auth.UserManager, hasher auth.Hasher, c config.CookieConfig) bool {
	// Get the requested user
	user, err := users.Get(auth.Fields{"name": username})
	if err != nil {
		// TODO Distinguish between errors and NotFound
		return false
	}

	// TODO hasher should be obtained from the password string
	if !auth.CheckPassword(hasher, password, user.Password()) {
		return false
	}

	// Create a new session
	session, err := sessions.Create(user)
	if err != nil {
		panic(fmt.Sprintf("server: could not create session: %s", err))
	}

	// Create a cookie
	auth.SetCookie(w, c, session)
	return true
}
