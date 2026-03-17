package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
)

func main() {
	jwksURL := "http://keycloak:8080/realms/fin_scale_realm/protocol/openid-connect/certs"

	k, err := keyfunc.NewDefault([]string{jwksURL})
	if err != nil {
		log.Fatalf("Failed to create keyfunc: %v", err)
	}

	http.HandleFunc("/validate", func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, k.Keyfunc)

		if err != nil || !token.Valid {
			fmt.Printf("Token invalid: %v\n", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			userID, _ := claims["sub"].(string)
			email, _ := claims["email"].(string)
			name, _ := claims["preferred_username"].(string)
			if !ok {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			w.Header().Set("X-User-ID", userID)
			w.Header().Set("X-User-Email", email)
			w.Header().Set("X-User-Name", name)
			log.Printf("User %s validated successfully", userID)
		}

		w.WriteHeader(http.StatusOK)
	})

	log.Println("Auth validator started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
