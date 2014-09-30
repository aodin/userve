package server

import (
	"log"
	"net/http"

	"github.com/aodin/volta/auth"
)

type PublicHandle func(http.ResponseWriter, *http.Request) ServerError

func (f PublicHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Recover from panic
	defer func() {
		if pnc := recover(); pnc != nil {
			// TODO Internal logging
			// TODO How to attach a 500 error template?
			log.Printf("panic in PublicHandle: %s\n", pnc)
			http.Error(w, "500", 500)
		}
	}()

	// Handle any other errors
	if err := f(w, r); err != nil {
		http.Error(w, err.Message(), err.Code())
	}
}

type Handle func(http.ResponseWriter, *http.Request, User) ServerError

// TODO Ugly, but we want templates and configuration
type Handler struct {
	f   Handle
	srv *Server
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Recover from panic
	defer func() {
		if pnc := recover(); pnc != nil {
			// TODO Internal logging
			// TODO How to attach a 500 error template?
			log.Printf("panic in Handle: %s\n", pnc)
			http.Error(w, "500", 500)
		}
	}()

	// Attempt to match the user, otherwise create an AnonUser
	var user User

	// Check if the session is valid
	// How to get server functions?
	cookie, err := r.Cookie(h.srv.config.Cookie.Name)
	if err == nil {
		u := auth.GetUserIfValidSession(h.srv.sessions, h.srv.users, cookie.Value)
		if u == nil || u.ID() == 0 {
			user = &AnonUser{}
		} else {
			// Convert the user returned from auth to a regular user
			// TODO awkward conversion - more parity?
			// TODO Other fields?
			user = &AuthUser{
				id:   u.ID(),
				name: u.Name(),
			}
		}
	} else {
		user = &AnonUser{}
	}

	// Handle any other errors
	if err := h.f(w, r, user); err != nil {
		http.Error(w, err.Message(), err.Code())
	}
}
