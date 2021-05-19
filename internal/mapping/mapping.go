package mapping

import (
	"books-api/internal"
	"books-api/pkg"

	"github.com/mmuoDev/commons/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

//generateHashPassword generates password hash from a string
func generateHashPassword(password string) (string, error) {
	bb, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", errors.Wrap(err, "Mapping- unable to hash password")
	}
	return string(bb), nil
}

//ToDBAuthor maps author request to internal author
func ToDBAuthor(r pkg.AuthorRequest) (internal.Author, error) {
	password, err := generateHashPassword(r.Password)
	if err != nil {
		return internal.Author{}, errors.Wrap(err, "Mapping - unable to hash password")
	}
	return internal.Author{
		ID:        uuid.GenV4(),
		Pseudonym: r.Pseudonym,
		Username:  r.Username,
		Password:  password,
	}, nil
}

//ToAuth maps token details to auth
func ToAuth(t *internal.Token, aID string) pkg.Auth {
	return pkg.Auth{
		AuthorID:    aID,
		AccessToken: t.Access,
	}
}
