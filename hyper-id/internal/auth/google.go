package auth

import (
	"context"

	"google.golang.org/api/idtoken"
)

const GoogleClientID = "306788201596-phllg7davd9ib2vkm1adrqrrv8r6ogoi.apps.googleusercontent.com"

func VerifyGoogleToken(tokenString string) (*idtoken.Payload, error) {
	// Ye function token ki authenticity check karta hai direct Google se
	payload, err := idtoken.Validate(context.Background(), tokenString, GoogleClientID)
	if err != nil {
		return nil, err
	}
	return payload, nil
}
