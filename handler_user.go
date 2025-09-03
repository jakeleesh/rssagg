package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jakeleesh/rssagg/internal/database"
)

// HTTP Handlers in Go, function signatures can't change
// But want to pass into function additional data
// So by making function a method, function signature remains the same, still just accepts 2 parameters
// But now have additional data stored on struct can gain access to
func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	// Handler needs to take as input a JSON body, expect parameters
	type parameters struct {
		Name string `json:"name"`
	}
	// Parse request body into struct
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	// Want to decode into an instance of parameter struct
	// Pointer into parameters
	err := decoder.Decode(&params)
	if err != nil {
		// Anything goes wrong, use use Handler function with error
		// Something goes wrong, probably client side so pass in 400
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		// Return because done if there is an issue
		return
	}

	// Use Database to create a new user
	// This method sqlc generated, accepts a context and CreateUserParams
	// r.Context() that's context for this request
	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		// user's name will be whatever was passed in HTTP request in the body
		Name: params.Name,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't create user: %v", err))
		return
	}

	// Rather than respond with database User, respond with our User
	// 201 is the created code
	respondWithJSON(w, 201, databaseUserToUser(user))
}

// New Handler for getting users
// This is an authenticated endpoint
// In order to create a user, don't need API key
// But if want get user info, have to give API key
func (apiCfg *apiConfig) handleGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJSON(w, 200, databaseUserToUser(user))
}
