package server

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/aodin/volta/config"
)

func TestServer(t *testing.T) {
	assert := assert.New(t)

	// Start a server a test the login / logout workflow
	// TODO simplify configuration
	c := config.DefaultConfig("")
	c.StaticDir = "../static"
	c.TemplateDir = "../templates"

	srv := New(c)

	// TODO Add these views in the creator function
	srv.AddPublicRoute("/", srv.Root)
	srv.AddPublicRoute("/login", srv.Login)
	srv.AddPublicRoute("/logout", srv.Logout)

	go srv.ListenAndServe()

	// Perform some simple requests
	response, err := http.Get("http://localhost" + c.Address())
	assert.Nil(err)

	// TODO Assert response function
	assert.Equal(response.StatusCode, 200)

	// TODO clean up test server? use httptest.NewServer?
}
