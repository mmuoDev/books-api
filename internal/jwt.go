package internal

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"errors"

	"github.com/dgrijalva/jwt-go"
	pkgErr "github.com/pkg/errors"
)

//TokenMetaData represents metadata of the token
type TokenMetaData struct {
	UserID     string
}

//GenerateJWT generates a token
func GenerateJWT(id string) (*Token, error) {
	expires := time.Now().Add(time.Minute * 15).Unix()

	//Access token
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = id
	claims["exp"] = expires
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(os.Getenv("JWT_ACCESS_SECRET")))
	if err != nil {
		return nil, pkgErr.Wrap(err, "jwt - unable to generate access token")
	}

	return &Token{Access: t}, nil
}

//verifyToken verifies signing mtd
func verifyToken(r *http.Request) (*jwt.Token, error) {
	ts, err := getToken(r)
	if err != nil {
		return nil, pkgErr.Wrap(err, "jwt - Token not found")
	}
	token, err := jwt.Parse(ts, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, pkgErr.Wrapf(err, "Invalid token!")
	}
	return token, nil
}

//getToken returns token from the header
func getToken(r *http.Request) (string, error) {
	bearerToken := r.Header.Get("Authorization")
	s := strings.Split(bearerToken, " ")
	if len(s) == 2 {
		return s[1], nil
	}
	return "", errors.New("jwt - Token not found")
}

//isTokenValid checks if token is valid
func isTokenValid(r *http.Request) bool {
	token, err := verifyToken(r)
	if err != nil {
		return false
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return false
	}
	return true
}

func GetTokenMetaData(r *http.Request) (*TokenMetaData, error) {
	token, err := verifyToken(r)
	if err != nil {
		return &TokenMetaData{}, pkgErr.Wrap(err, "jwt - unable to retrieve token")
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["user_id"].(string)
		if !ok {
			return nil, fmt.Errorf("userID metadata not found in token")
		}
		return &TokenMetaData{
			UserID:     userID,
		}, nil
	}
	return &TokenMetaData{}, err
}
