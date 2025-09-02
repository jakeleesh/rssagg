package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// Server we're building going to be JSON REST API
// Means all request bodies coming in and going back will have JSON format
// Helper function make it easier send JSON responses

// Responding with arbitrary error messages
// Instead of taking a payload, take a message string
// Function basically format message into a consistent JSON object every time
func respondWithError(w http.ResponseWriter, code int, msg string) {
	// Error codes in 400 range are client side errors, don't need to know about them.
	// Means using our API in weird way.
	// Need to know 500 level error code because means bug on our end
	if code > 499 {
		log.Println("Responding with 5XX error:", msg)
	}
	// Responding with specific structure of JSON
	// Take struct and add JSON tags to specify how we want to unmarshal, convert struct into JSON object
	type errResponse struct {
		// Struct has 1 field, Error
		// Add this JSON tag to say this key should marshal to error
		// Saying I have error field, want key for field to be error
		Error string `json:"error`
	}

	respondWithJSON(w, code, errResponse{
		Error: msg,
	})
}

// Helper function for responding with arbitrary JSON
// Takes:
// ResponseWriter HTTP handlers use,
// status code to respond with,
// interface which is a JSON structure
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	// Marshal payload into JSON string
	// Return as bytes so that can write it in binary format directly to HTTP
	dat, err := json.Marshal(payload)
	// If fails
	if err != nil {
		// Log it and print what we tried to marshal
		log.Printf("Failed to marshal JSON response: %v", payload)
		// Write a Header to response
		// Use 500, say something went wrong on our end, internal error
		w.WriteHeader(500)
		return
	}
	// Add a Header to response
	// Say we're responding with JSON
	// key: Content-Type, value: application/json
	// Adds a response header to HTTP request saying responding with content type of application/json
	// Standard value for JSON response
	w.Header().Add("Content-Type", "application/json")
	// Use the passed in response code
	w.WriteHeader(code)
	// Write data itself, pass in JSON data
	w.Write(dat)
}
