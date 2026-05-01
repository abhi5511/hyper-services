package main

import (
	"fmt"
	"log"
	"net/http"
	"os" // Environment variables read karne ke liye

	"hyper-id/internal/api"
	"hyper-id/internal/auth"
	"hyper-id/internal/db"

	"github.com/joho/godotenv" // .env file load karne ke liye
)

// Middleware to enable CORS for development
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Google-Token, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	// 1. Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Note: .env file not found, using system environment variables")
	}

	fmt.Println("Starting Hyper ID Service...")

	// 2. Load RSA Private Key
	err = auth.InitKeys("certs/private.pem")
	if err != nil {
		log.Fatalf("Critical: Could not load RSA keys: %v", err)
	}

	// 3. Initialize Database Connection using Environment Variables
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Database Connect
	err = db.InitDB(dbHost, dbPort, dbUser, dbPass, dbName)
	if err != nil {
		log.Fatalf("Critical: Database connection failed: %v", err)
	}
	fmt.Println("Database Connected Successfully!")

	// 4. Setup Routes
	mux := http.NewServeMux()

	// API Routes
	mux.HandleFunc("/auth/google", api.HandleGoogleLogin)
	mux.HandleFunc("/auth/onboarding", api.RequireAuth(api.HandleSetUsername))

	// Static Files
	fs := http.FileServer(http.Dir("frontend"))
	mux.Handle("/", fs)

	// 5. Start Server
	port := ":8080"
	fmt.Printf("Hyper ID Auth Server running on http://localhost%s\n", port)

	if err := http.ListenAndServe(port, enableCORS(mux)); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
