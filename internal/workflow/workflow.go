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
type RetrieveBooksFunc func(params pkg.QueryParams) ([]pkg.Book, error)

//DeleteBookByIDFunc deletes a book by id
type DeleteBookByIDFunc func(aID, bID string) error

//RetrieveBookByIDFunc retrieves book by id
type RetrieveBookByIDFunc func(bID string) (pkg.Book, error)

//UpdateBookFunc updates a book
type UpdateBookFunc func(bID string, r pkg.BookUpdateRequest) error

//UpdateBook updates a book
func UpdateBook(updateBook db.UpdateBookFunc) UpdateBookFunc {
	return func(bID string, r pkg.BookUpdateRequest) error {
		changes := mapping.ToDBBookUpdated(r)
		if err := updateBook(bID, changes); err != nil {
			return pkgErr.Wrapf(err, "workflow - unable to update book with id=%s", bID)
		}
		return nil
	}
}

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
func RetrieveBooks(retrieve db.RetrieveBooksFunc, retrieveAuthor db.RetrieveAuthorByIDFunc) RetrieveBooksFunc {
	return func(params pkg.QueryParams) ([]pkg.Book, error) {
		books, err := retrieve(params)
		if err != nil {
			return []pkg.Book{}, pkgErr.Wrap(err, "Workflow - error retrieving books")
		}
		return mapping.ToDTOBooks(books, retrieveAuthor), nil
	}
}

//DeleteBookByID deletes book by id for authenticated user
func DeleteBookByID(delete db.DeleteBookByIDFunc) DeleteBookByIDFunc {
	return func(aID, bID string) error {
		if err := delete(aID, bID); err != nil {
			return pkgErr.Wrapf(err, "workflow - error deleting book with id=%s by authorID=%s", bID, aID)
		}
		return nil
	}
}

//RetrieveBookByID retrieves book by id
func RetrieveBookByID(retrieve db.RetrieveBookByIDFunc, retrieveAuthor db.RetrieveAuthorByIDFunc) RetrieveBookByIDFunc {
	return func(bID string) (pkg.Book, error) {
		b, err := retrieve(bID)
		if err != nil {
			return pkg.Book{}, pkgErr.Wrapf(err, "Workflow - error retrieving book with id=%s", bID)
		}
		a, _ := retrieveAuthor(b.AuthorID)
		return mapping.ToDTOBook(b, a.Pseudonym), nil
	}
}

//validatePassword validates password for a user
func validatePassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
