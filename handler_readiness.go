package main

import "net/http"

// Very specific function signature
// Have to use if want to define HTTP Handler the way Go standard library expects
// Always takes ResponseWriter as first parameter
// Pointer to HTTP request as 2nd parameter
func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	// Call respondWithJSON function
	// Pass in ResponseWriter,
	// Want to respond with 200,
	// Response payload. In this case, all we care about is 200 so respond with empty struct, should marshal to empty JSON object
	respondWithJSON(w, 200, struct{}{})
}
