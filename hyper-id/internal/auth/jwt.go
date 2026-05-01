package auth

import (
	"crypto/rsa"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// privateKey ko hum global rakhenge taaki GenerateHyperToken ise use kar sake
var privateKey *rsa.PrivateKey

type HyperClaims struct {
	HID      string `json:"hid"`
	Username string `json:"username"`
	Status   string `json:"status"`
	jwt.RegisteredClaims
}

// InitKeys function ko main.go se call karna hoga server start hote waqt
func InitKeys(privateKeyPath string) error {
	// 1. .pem file read karo
	keyBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return err
	}

	// 2. PEM data ko RSA Private Key mein parse karo
	key, err := jwt.ParseRSAPrivateKeyFromPEM(keyBytes)
	if err != nil {
		return err
	}

	privateKey = key
	return nil
}

func GenerateHyperToken(hid, username, status string) (string, error) {
	claims := HyperClaims{
		HID:      hid,
		Username: username,
		Status:   status,
		RegisteredClaims: jwt.RegisteredClaims{
			// Token 24 ghante tak valid rahega
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "hyper-id",
		},
	}

	// RS256 use kar rahe hain (Asymmetric Signing)
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	
	// Private Key se sign karke string return karo
	return token.SignedString(privateKey)
}

func VerifyHyperToken(tokenString string) (*HyperClaims, error) {
	claims := &HyperClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// RS256 ko verify karne ke liye Public Key chahiye hoti hai
		return &privateKey.PublicKey, nil 
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	return claims, nil
}