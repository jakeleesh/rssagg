package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	// Environment doesn't exist in current shell session
	// Use package to grab environment variables
	godotenv.Load(".env")

	// Read PORT variable by key
	portString := os.Getenv("PORT")
	if portString == "" {
		// log.Fatal will exit program immediately with Error Code 1 and message
		log.Fatal("PORT is not found in the environment")
	}

	// Spin up Server
	// New Router Object
	router := chi.NewRouter()

	// cors configuration from cors package installed
	// Essentially telling Server to send extra HTTP Headers, tell browsers allow to use these
	router.Use(
		cors.Handler(
			cors.Options{
				// Allow send requests to http or https
				AllowedOrigins: []string{"https://*", "http://*"},
				// Allow methods
				AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
				// Allow send any Headers
				AllowedHeaders:   []string{"*"},
				ExposedHeaders:   []string{"Link"},
				AllowCredentials: false,
				MaxAge:           300,
			},
		),
	)

	// Create new Router
	v1Router := chi.NewRouter()
	// Hook up HTTP Handler to a specific HTTP method and path
	// Handle /healthz path with handlerReadiness function
	// Name healthz, kubernetes standard to see if server is live and running
	// POST request get 200 not intention
	// healthz endpoint should only be accessible by GET request
	// Rather than using v1Router.handleFunc, use v1Router.Get. Scope hanlder to only fire on GET requests.
	v1Router.Get("/healthz", handlerReadiness)
	// Hook up error handler
	v1Router.Get("/err", handlerErr)

	// Create v1Router is because going to mount
	// Nesting v1Router under /v1 path
	// Full path for request will be: /v1/healthz
	// So that if make changes in future, can have 2 handlers, v1 and v2 for API. Standard practie.
	router.Mount("/v1", v1Router)

	// Connect Router to HTTP Server
	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Server starting on port %v", portString)
	// Returns an Error
	// ListenAndServe will block, just stop and starts handling HTTP Requests
	// Nothing SHOULD be returned, Server should run forever
	err := srv.ListenAndServe()
	// Anything goes wrong in process of handling requests, error returned
	if err != nil {
		// Log and exit program
		log.Fatal(err)
	}

	fmt.Println("Port:", portString)
}
