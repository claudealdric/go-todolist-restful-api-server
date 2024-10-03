package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("secret_key") // TODO: change

type LoginCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"accessToken"`
}

func (s *Server) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var credentials LoginCredentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, "Error parsing the JSON body", http.StatusBadRequest)
		return
	}

	if !s.store.ValidateUserCredentials(
		credentials.Email,
		credentials.Password,
	) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{"exp": expirationTime},
	)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Println("error signing the JWT:", err)
		http.Error(w, "Error creating the JWT", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	response := LoginResponse{AccessToken: tokenString}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("error encoding response: %v", err)
	}

}
