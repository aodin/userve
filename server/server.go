package server

import (
	"crypto/sha1"
	"fmt"
	"log"
	"net/http"

	"github.com/aodin/aspect"
	"github.com/aodin/volta/auth"
	"github.com/aodin/volta/config"
	"github.com/aodin/volta/templates"
)

type Server struct {
	attrs     Attrs
	config    config.Config
	db        *aspect.DB
	templates templates.Templates
	users     auth.UserManager
	sessions  auth.SessionManager
	hasher    auth.Hasher
}

func (srv *Server) Attrs(b map[string]interface{}) map[string]interface{} {
	srv.attrs.Merge(b)
	return b
}

func (srv *Server) Execute(w http.ResponseWriter, name string, data Attrs) {
	if data == nil {
		data = Attrs{}
	}
	srv.attrs.Merge(data)
	if err := srv.templates.Execute(name, "", w, data); err != nil {
		panic(fmt.Sprintf("server: could execute template: %s", err))
	}
}

func (srv *Server) ExecuteLayout(w http.ResponseWriter, name string, d Attrs) {
	if d == nil {
		d = Attrs{}
	}
	srv.attrs.Merge(d)
	if err := srv.templates.Execute(name, "layout", w, d); err != nil {
		panic(fmt.Sprintf("server: could not execute layout: %s", err))
	}
}

func (srv *Server) AddPublicRoute(route string, h PublicHandle) {
	http.Handle(route, h)
}

func (srv *Server) AddUserRoute(route string, h Handle) {
	http.Handle(route, h)
}

func (srv *Server) ListenAndServe() error {
	log.Printf("Starting server on %s\n", srv.config.Address())
	return http.ListenAndServe(srv.config.Address(), nil)
}

func New(c config.Config) (srv *Server) {
	srv = &Server{config: c, attrs: Attrs{"StaticURL": c.StaticURL}}

	// Parse templates
	var err error
	srv.templates, err = templates.Parse(c.TemplateDir, "layout.html")
	if err != nil {
		log.Fatalf("error parsing templates: %s", err)
	}

	// Serve static files
	files := http.FileServer(http.Dir(c.StaticDir))
	http.Handle(c.StaticURL, http.StripPrefix(c.StaticURL, files))

	// Create the auth
	srv.hasher = auth.NewPBKDF2Hasher("test", 1, sha1.New)
	srv.users = auth.UsersInMemory(srv.hasher)

	// Create a new user
	_, err = srv.users.Create("admin", "admin", auth.Fields{"isAdmin": true})
	if err != nil {
		log.Fatalf("error creating admin: %s", err)
	}

	// Create a new sessions manager
	srv.sessions = auth.SessionsInMemory(c.Cookie)

	return
}

// Handlers
// ========

func (srv *Server) Root(w http.ResponseWriter, r *http.Request) ServerError {
	srv.ExecuteLayout(w, "root.html", nil)
	return nil
}

func (srv *Server) Login(w http.ResponseWriter, r *http.Request) ServerError {
	// Unless it is a POST, just display the login form
	if r.Method != "POST" {
		srv.ExecuteLayout(w, "login.html", nil)
		return nil
	}

	username := r.FormValue("username")

	// TODO invalid username / password errors
	if AuthLogin(
		w,
		username,
		r.FormValue("password"),
		srv.sessions,
		srv.users,
		srv.hasher,
		srv.config.Cookie,
	) {
		http.Redirect(w, r, "/", 302)
		return nil
	}

	data := map[string]interface{}{
		"Username": username,
		"Message":  "Invalid credentials",
	}

	srv.ExecuteLayout(w, "login.html", data)
	return nil
}

func (srv *Server) Logout(w http.ResponseWriter, r *http.Request) ServerError {
	// Remove the session
	cookie, err := r.Cookie(srv.config.Cookie.Name)
	if err != nil {
		panic(err)
	}

	if err = srv.sessions.Delete(cookie.Value); err != nil {
		panic(fmt.Sprintf("server: could not delete session cookie: %s", err))
	}

	http.Redirect(w, r, "/", 302)
	return nil
}
