package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/iboutilierRED/backend-workshop/ch-3/wikievents/internal/models"
)

func SignUpHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// only accept POST methods for sign up
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed. Use POST", http.StatusMethodNotAllowed)
			return
		}

		var user models.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}
		// save user and return id
		w.Write([]byte("User created! You may now login with your email and password."))
	}
}

func LoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed. Use POST", http.StatusMethodNotAllowed)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")

		// stats := sc.GetStats()
		// json.NewEncoder(w).Encode(stats)
	}
}
