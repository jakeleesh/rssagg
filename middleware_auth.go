package main

import (
	"fmt"
	"net/http"

	"github.com/jakeleesh/rssagg/internal/auth"
	"github.com/jakeleesh/rssagg/internal/database"
)

// Create bew Handler that will allow users to create a new feed
// That Handler going to need same logic we have in GetUser Handler
// Rather than copying code into every Handler that's authenticated, build middleware to DRY code

// Define new custom Handler
// Looks like almost like regular HTTP Handler
// Only difference is that it includes a 3rd parameter, has User
// Makes sense, 3rd one is Authenticated User
type authedHandler func(http.ResponseWriter, *http.Request, database.User)

// Problem with this Handler type is that it doesn't match function signature of an HTTP Handler
// Functions with just ResponseWriter and request
// Create new function
// Method on apiConfig so that is has access to database
// Job is to taje authedHandler and return a HandlerFunc so that can use with chi router
func (apiCfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	// Return a Closure, anonymous function
	// Same function signature as HTTP Handler
	// Only difference is, have access to everything withing apiConfig, able to query database
	return func(w http.ResponseWriter, r *http.Request) {
		// Rip out code from GetUser Handler
		// Get API Key from request
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			// Error respond with 403 for creating a user
			respondWithError(w, 403, fmt.Sprintf("Auth error: %v", err))
			return
		}

		// Have API Key, can use database query
		// Context package in standard library
		// Basically gives a way to track something that's happening accross multiple goroutines
		// Most important thing you can do with context is cancel it
		// Cancelling context would effectively kill HTTP request
		// Make sure use current context
		// Every HTTP request has a context on it
		// Should use that context in any calls make within the handler that requires context in case cancellations happen
		// Grab User using API Key
		user, err := apiCfg.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Couldn't get user: %v", err))
			return
		}

		// By the time get to calling the Handler, able to give actual user from database
		handler(w, r, user)
	}
}
