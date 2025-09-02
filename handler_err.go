package main

import "net/http"

func handlerErr(w http.ResponseWriter, r *http.Request) {
	// Instead of passing in an empty struct, say
	// 400 status code client error
	respondWithError(w, 400, "Something went wrong")
}
