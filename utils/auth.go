package utils

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"rc-h3-webapp/models"
	"time"
)

var SecretKey string

func hashMD5(input string) string {
	hasher := md5.New()                        // Create a new MD5 hash
	hasher.Write([]byte(input))                // Write the input string to the hasher
	return hex.EncodeToString(hasher.Sum(nil)) // Return the hex-encoded hash
}

func CheckPassword(inputPassword string, storedHash string) bool {
	// Hash the input password
	hashedPassword := hashMD5(inputPassword)

	// Compare the hashed password with the stored hash
	return hashedPassword == storedHash
}

func GenerateJWT(username string, userId int) (string, error) {
	// Set token expiration time
	expirationTime := time.Now().Add(2 * time.Hour) // Token valid for 1 hour

	// Create the claims
	claims := &models.JwtClaims{
		Username: username,
		UserId:   userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(), // Set expiration time
		},
	}

	// Generate the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with your secret key
	secretKey := []byte(SecretKey) // Replace with your actual secret key
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GenerateSecretKey(length int) {
	// Create a byte slice to hold the random bytes
	key := make([]byte, length)

	// Read random bytes into the slice
	if _, err := rand.Read(key); err != nil {
		panic(err)
	}

	// Encode the bytes to a Base64 string for easy usage
	SecretKey = base64.StdEncoding.EncodeToString(key)
}

func VerifyJWT(tokenString string) (*models.JwtClaims, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &models.JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Ensure that the token method is valid
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Return the secret key for verification
		return []byte(SecretKey), nil
	})

	if err != nil {
		return nil, err // Return the error if token parsing fails
	}

	// Check if the token is valid and claims are available
	if claims, ok := token.Claims.(*models.JwtClaims); ok && token.Valid {
		return claims, nil // Return the claims if the token is valid
	}

	return nil, fmt.Errorf("invalid token") // Return an error if the token is invalid
}
