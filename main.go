package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/jakeleesh/rssagg/internal/database"
	"github.com/joho/godotenv"

	// Underscore to say include this code in program even though not calling it directly
	_ "github.com/lib/pq"
)

// Use Database in code
// struct hold connection to database
type apiConfig struct {
	// Exposed by code generated using sqlc
	DB *database.Queries
}

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

	// Import database connection
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not found in environment")
	}

	// Connect to database
	// Go standard library has built-in SQL package
	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Can't connect to database:", err)
	}

	db := database.New(conn)
	// New API Config
	// Can pass into our handlers so that they have access to database
	apiCfg := apiConfig{
		// Takes in database.queries
		// Have sql.db so need to convert into a connection
		DB: db,
	}

	// Hook up startScraping to main function
	// Call before ListenAndServe() because server blocks and waits forever for incoming requests
	// Call it on a new goroutine so doesn't interrupt main
	// because startScraping is never going to return, it's long running functio, infinite for loop
	go startScraping(db, 10, time.Minute)

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
	// Hook up createUser Handler
	// Be POST Request
	v1Router.Post("/users", apiCfg.handlerCreateUser)
	// Hook up GetUser Handler to GET HTTP method
	// Same path, different method
	// Call middlewareAuth to convert GetUser Handler into standard HTTP Handler
	// Calling middlewareAuth to get authenticated user and then calling back the GetUser Handler
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handleGetUser))

	// Creating a resouce, use POST
	v1Router.Post("/feeds", apiCfg.middlewareAuth((apiCfg.handlerCreateFeed)))
	v1Router.Get("/feeds", apiCfg.handlerGetFeeds)

	v1Router.Post("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollow))
	v1Router.Get("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerGetFeedFollows))
	// Authenticated
	// Need feedFollowID and DELETE request
	// HTTP DELETE request don't typically have body
	// More conventional to pass ID in path
	v1Router.Delete("/feed_follows/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.handlerDeleteFeedFollow))

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
	err = srv.ListenAndServe()
	// Anything goes wrong in process of handling requests, error returned
	if err != nil {
		// Log and exit program
		log.Fatal(err)
	}

	fmt.Println("Port:", portString)
}
