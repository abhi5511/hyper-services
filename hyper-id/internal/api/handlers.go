package api

import (
	"encoding/json"
	"net/http"

	"hyper-id/internal/auth"
	"hyper-id/internal/db"
)

// Response structure for better type safety
type LoginResponse struct {
	Token  string `json:"token"`
	Status string `json:"status"` // 'pending_onboarding' ya 'active'
}

func HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	// 1. Google Token Header se nikalna
	idToken := r.Header.Get("X-Google-Token")
	if idToken == "" {
		http.Error(w, "Missing Google Token", http.StatusBadRequest)
		return
	}

	// 2. Google Token Verify karo
	payload, err := auth.VerifyGoogleToken(idToken)
	if err != nil {
		http.Error(w, "Invalid Google Token", http.StatusUnauthorized)
		return
	}

	// 3. Email extract karna
	email, ok := payload.Claims["email"].(string)
	if !ok {
		http.Error(w, "Email not provided by Google", http.StatusBadRequest)
		return
	}

	// 4. Database mein User check karo ya create karo
	user, err := db.FindOrCreateUser(payload.Subject, email)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// 5. Hyper ID JWT Generate karo
	token, err := auth.GenerateHyperToken(user.HID, user.Username, user.AccountStatus)
	if err != nil {
		http.Error(w, "Token generation failed", http.StatusInternalServerError)
		return
	}

	// --- NEW: Multi-Project Redirect Logic ---

	// Request URL se parameters pakadna (e.g., ?client_id=...&redirect_uri=...)
	clientID := r.URL.Query().Get("client_id")
	redirectURI := r.URL.Query().Get("redirect_uri")

	w.Header().Set("Content-Type", "application/json")

	// Base Response
	response := map[string]interface{}{
		"token":  token,
		"status": user.AccountStatus,
	}

	// Agar external app se request aayi hai, toh redirect_url build karo
	if clientID != "" && redirectURI != "" {
		// Yahan user ko uske original app par wapas bhejne ka path tayaar hai
		response["redirect_url"] = redirectURI + "?token=" + token + "&status=" + user.AccountStatus
	}

	// 6. Final Response bhejo
	json.NewEncoder(w).Encode(response)
}
