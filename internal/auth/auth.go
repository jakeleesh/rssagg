// Abstract the logic
// Create a new package
package auth

import (
	"errors"
	"net/http"
	"strings"
)

// GetAPIKey extracts an API Key from the headers of an HTTP request
// Go into headers and see if can find the API Key, otherwise return an error
// Looking for header of this format, Example:
// Authorization: ApiKey {apikey}
func GetAPIKey(headers http.Header) (string, error) {
	// Using http standard library
	// See value associated with this header key
	val := headers.Get("Authorization")
	if val == "" {
		return "", errors.New("no authentication info found")
	}

	vals := strings.Split(val, " ")
	// Expecting value of key is 2 specific values seperated by a space
	// ApiKey {apikey} ("ApiKey" and THE ACTUAL apikey)
	if len(vals) != 2 {
		return "", errors.New("malformed auth header")
	}
	if vals[0] != "ApiKey" {
		return "", errors.New("malformed first part of auth header")
	}
	return vals[1], nil
}
