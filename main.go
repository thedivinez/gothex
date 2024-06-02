package main

import (
	"encoding/json"
	"log"

	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"github.com/thedivinez/gothex/api"
	"github.com/thedivinez/gothex/frontend"
	"github.com/thedivinez/gothex/pages"
	"github.com/thedivinez/gothex/router"
)

type Configs struct {
	GOOGLE_KEY      string `json:"GOOGLE_KEY"`
	GOOGLE_SECRET   string `json:"GOOGLE_SECRET"`
	GOOGLE_CALLBACK string `json:"GOOGLE_CALLBACK"`
}

type Server struct {
	configs *Configs
	pages   *frontend.Pages
	auth    *api.AuthHandler
}

func NewServer() *Server {
	data, err := godotenv.Read()
	if err != nil {
		log.Fatal(err)
	}
	resultBytes, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	var configs Configs
	if err := json.Unmarshal(resultBytes, &configs); err != nil {
		log.Fatal(err)
	}
	return &Server{
		configs: &configs,
		pages:   frontend.NewPages(),
		auth:    api.NewAuthHandler(),
	}
}

func main() {
	server := NewServer()
	gothex := router.NewGothexRouter().WithCustomErrorPageContent(
		pages.ErrorPageContent{
			Code:      404,
			Title:     "Error (404)",
			ErrorType: "Resource not found",
			Message:   "The requested resource could not be found but may be available again in the future.",
		},
		pages.ErrorPageContent{
			Code:      500,
			Title:     "Error (500)",
			ErrorType: "Internal Server Error",
			Message:   "An unexpected condition was encountered and prevented the request from succeeding.",
		},
	)

	gothic.Store = gothex.CookieStore
	goth.UseProviders(google.New(server.configs.GOOGLE_KEY, server.configs.GOOGLE_SECRET, server.configs.GOOGLE_CALLBACK))

	gothex.GET("/signin", server.pages.HandleSignIn)
	gothex.GET("/signup", server.pages.HandleSignUp)
	gothex.GET("/", server.pages.HandleHome, gothex.Protected)

	api := gothex.Group("/api")
	api.Any("/signin", server.auth.HandleSignIn)
	api.POST("/signup", server.auth.HandleSignUp)
	api.GET("/signin/callback", server.auth.HandleProvidersCallback)
	api.POST("/signout", server.auth.HandleSignOut, gothex.Protected)
	api.GET("/test-session", server.auth.HandleAuthError, gothex.Protected)

	if err := gothex.Run(); err != nil {
		log.Fatal("failed to start server", err)
	}
}
