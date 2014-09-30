package server

import (
	"log"
	"net/http"
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

func (f Handle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Recover from panic
	defer func() {
		if pnc := recover(); pnc != nil {
			// TODO Internal logging
			// TODO How to attach a 500 error template?
			log.Printf("panic in Handle: %s\n", pnc)
			http.Error(w, "500", 500)
		}
	}()

	// TODO Attempt to match the user, otherwise create an AnonUser
	anon := &AnonUser{}

	// Handle any other errors
	if err := f(w, r, anon); err != nil {
		http.Error(w, err.Message(), err.Code())
	}
}
