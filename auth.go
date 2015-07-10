package main

import "net/http"

func AuthenticateRequest(w http.ResponseWriter, r *http.Request) {
	//redirect user to login if they're unauthenticated
	authenticated := false
	if !authenticated {
		http.Redirect(w, r, "/register", http.StatusFound)
	}
}
