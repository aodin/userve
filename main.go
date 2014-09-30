package main

import (
	"log"

	"github.com/aodin/userve/server"
	"github.com/aodin/volta/config"
)

func main() {
	c := config.DefaultConfig("")
	c.StaticDir = "./static"
	c.TemplateDir = "./templates"

	srv := server.New(c)

	// TODO Add these views in the creator function
	srv.AddUserRoute("/", srv.Root)
	srv.AddPublicRoute("/login", srv.Login)
	srv.AddPublicRoute("/logout", srv.Logout)

	log.Fatal(srv.ListenAndServe())
}
