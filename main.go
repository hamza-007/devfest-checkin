package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"fmt"

	godotenv "github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load("./.env"); err != nil {
		fmt.Println("error loading env file" ,err)
		return
	}

	mux := http.NewServeMux()

	handler := CheckHeader("X-API-ACCESS", os.Getenv("X_API_ACCESS"))(mux)
	
	// QR code endpoint
	mux.HandleFunc("/qrcode", func(w http.ResponseWriter, r *http.Request) {
		// Get "text" query parameter
		userID := r.URL.Query().Get("user_id")
		if userID == "" {
			http.Error(w, "missing 'user_id' query param", http.StatusBadRequest)
			return
		}

		if userID == "123" {
			resp := SuccessResponse{
				Status:  "success",
				Message: "User registered !",
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
		} else {
			resp := SuccessResponse{
				Status:  "failed",
				Message: "User not found !",
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
		}

	})


	log.Fatal(http.ListenAndServe(":8080", handler))
}

func CheckHeader(headerName, expectedValue string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get(headerName) != expectedValue {
				http.Error(w, "Forbidden: invalid access token", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r) // header OK, call next handler
		})
	}
}

type SuccessResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
