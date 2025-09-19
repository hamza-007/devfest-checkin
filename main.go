package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	mail "devfest-checkin/utils/mail"
	qr "devfest-checkin/utils/qr"

	firestore "cloud.google.com/go/firestore"
	godotenv "github.com/joho/godotenv"
	option "google.golang.org/api/option"
)

func main() {

	ctx := context.Background()
	sa := option.WithCredentialsFile("serviceAccountKey.json")

	client, err := firestore.NewClient(ctx, "devfest-ac5cd", sa)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	if err := godotenv.Load("./.env"); err != nil {
		fmt.Println("error loading env file", err)
		return
	}

	collection := client.Collection("users")

	// Add a document
	doc, _, err := collection.Add(ctx, map[string]interface{}{
		"name":  "Hamza",
		"email": "hamza@example.com",
		"age":   25,
	})
	if err != nil {
		log.Fatalf("Failed adding document: %v", err)
	}
	fmt.Println("Added document with ID:", doc.ID)

	return 

	mux := http.NewServeMux()

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
			http.Error(w, "invalid user", http.StatusBadRequest)
			return
		}

	})

	handler := CheckHeader("X-API-ACCESS", os.Getenv("X_API_ACCESS"))(mux)

	png, err := qr.GenerateQR("http://127.0.0.1:8023?user_id=123")
	if err != nil {
		return
	}

	sender := mail.New("hamzaabdellaoui26648999@gmail.com")
	err = sender.SendMail(png)
	if err != nil {
		fmt.Println("error sending email", err)
	}
	fmt.Println("email sent ! ")
	log.Fatal(http.ListenAndServe(":8023", handler))
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
