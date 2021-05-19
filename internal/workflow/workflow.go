package workflow

import (
	"books-api/internal"
	"books-api/internal/db"
	"books-api/internal/mapping"
	"books-api/pkg"

	"errors"

	pkgErr "github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

//AddAuthorFunc returns functionality to add an author
type AddAuthorFunc func(r pkg.AuthorRequest) error

//AuthenticateFunc authenticates a user
type AuthenticateFunc func(r pkg.AuthRequest) (pkg.Auth, error)

//AddBookFunc adds a book for an author
type AddBookFunc func(r pkg.BookRequest, aID string) error 

//RetrieveBooksFunc retrieves books
type RetrieveBooksFunc func() ([]pkg.BookRequest, error)

//DeleteBookByIDFunc deletes a book by id
type DeleteBookByIDFunc func(aID, bID string) error

//AddAuthor adds an author
func AddAuthor(addAuthor db.AddAuthorFunc) AddAuthorFunc {
	return func(r pkg.AuthorRequest) error {
		a, err := mapping.ToDBAuthor(r)

		if err != nil {
			return pkgErr.Wrap(err, "Workflow - unable to map internal user to db")
		}
		if err := addAuthor(a); err != nil {
			return pkgErr.Wrap(err, "Workflow - error adding new user")
		}
		return nil
	}
}

//Authenticate authenticates a user and generates a token
func Authenticate(retriveAuthor db.RetrieveAuthorByUsernameFunc) AuthenticateFunc {
	return func(r pkg.AuthRequest) (pkg.Auth, error) {
		retrieve, err := retriveAuthor(r.Username)
		if err != nil {
			return pkg.Auth{}, pkgErr.Wrapf(err, "Workflow - No author found for username=%s", r.Username)
		}
		if vp := validatePassword(retrieve.Password, r.Password); !vp {
			return pkg.Auth{}, errors.New("Workflow - Incorrect auth credentials")
		}
		aID := retrieve.ID.Val()
		t, err := internal.GenerateJWT(aID)
		if err != nil {
			return pkg.Auth{}, pkgErr.Wrap(err, "Workflow - Unable to generate tokens")
		}
		return mapping.ToAuth(t, aID), nil
	}
}

//AddBook adds a book
func AddBook(addBook db.AddBookFunc) AddBookFunc {
	return func(r pkg.BookRequest, aID string) error {
		b := mapping.ToDBBook(r, aID)
		if err := addBook(b); err != nil {
			return pkgErr.Wrap(err, "Workflow - error adding new book")
		}
		return nil
	}
}

//RetrieveBooks retrieves books 
func RetrieveBooks(retrieve db.RetrieveBooksFunc) RetrieveBooksFunc {
	return func() ([]pkg.BookRequest, error) {
		books, err := retrieve()
		if err != nil {
			return []pkg.BookRequest{}, pkgErr.Wrap(err, "Workflow - error retrieving books")
		}
		return mapping.ToDTOBooks(books), nil
	}
}

//DeleteBookByID deletes book by id for authenticated user
func DeleteBookByID (delete db.DeleteBookByIDFunc) DeleteBookByIDFunc {
	return func(aID, bID string) error {
		if err := delete(aID, bID); err != nil {
			return pkgErr.Wrapf(err, "workflow - error deleting book with id=%s by authorID=%s", bID, aID)
		}
		return nil 
	}
}

//validatePassword validates password for a user
func validatePassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
