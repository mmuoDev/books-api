package internal

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

//GenerateJWT generates a token
func GenerateJWT(id string) (*Token, error) {
	//td := &internal.Token{}
	expires := time.Now().Add(time.Minute * 15).Unix()

	//Access token
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = id
	claims["exp"] = expires
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(os.Getenv("JWT_ACCESS_SECRET")))
	if err != nil {
		return nil, errors.Wrap(err, "Workflow - unable to generate access token")
	}

	return &Token{Access: t}, nil
}
