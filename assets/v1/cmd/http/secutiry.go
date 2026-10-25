package main

import "net/http"

// Helper for authentication check (Header or Cookie)
func isAuthorizedRequest(r *http.Request) bool {
	headerToken := r.Header.Get("X-SecretAccessToken")
	cookieToken := ""
	if cookie, err := r.Cookie("accessToken"); err == nil {
		cookieToken = cookie.Value
	}
	return (headerToken != "" && headerToken == serviceToken) || (cookieToken != "" && cookieToken == serviceToken)
}
