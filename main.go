package main

import (
	"log"

	"github.com/aodin/volta/config"
	"github.com/aodin/userve/server"
)

func main() {
	c := config.DefaultConfig("")
	c.StaticDir = "./static"
	c.TemplateDir = "./templates"

	srv := server.New(c)

	// TODO Add these views in the creator function
	srv.AddPublicRoute("/", srv.Root)
	srv.AddPublicRoute("/login", srv.Login)
	srv.AddPublicRoute("/logout", srv.Logout)

	log.Fatal(srv.ListenAndServe())
}
