package db

import (
	"books-api/internal"

	"github.com/mmuoDev/commons/mongo"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	authorsCollection = "authors"
	booksCollection   = "books"
)

//AddAuthorFunc returns functionality to add an author to DB
type AddAuthorFunc func(internal.Author) error

//RetrieveAuthorByIDFunc retrieves an author by author id
type RetrieveAuthorByIDFunc func(aid string) (internal.Author, error)

//RetrieveAuthorByUsernameFunc retrieves an author by username
type RetrieveAuthorByUsernameFunc func(username string) (internal.Author, error)

//AddAuthor adds an author to DB
func AddAuthor(dbProvider mongo.DbProviderFunc) AddAuthorFunc {
	return func(a internal.Author) error {
		col := mongo.NewCollection(dbProvider, authorsCollection)
		_, err := col.Insert(a)
		if err != nil {
			return errors.Wrap(err, "db - failure inserting an author")
		}
		return nil
	}
}

//RetrieveAuthorByID retrieves author by the author's id
func RetrieveAuthorByID(dbProvider mongo.DbProviderFunc) RetrieveAuthorByIDFunc {
	return func(aid string) (internal.Author, error) {
		col := mongo.NewCollection(dbProvider, authorsCollection)
		var a internal.Author
		if err := col.FindByID(aid, &a); err != nil {
			return a, errors.Wrap(err, "db - failure retrieving author")
		}
		return a, nil
	}
}

//RetrieveAuthorByUsername retrieves author by user
func RetrieveAuthorByUsername(dbProvider mongo.DbProviderFunc) RetrieveAuthorByUsernameFunc {
	return func(username string) (internal.Author, error) {
		col := mongo.NewCollection(dbProvider, authorsCollection)
		filter := bson.D{{"username", username}}
		var a internal.Author

		if err := col.FindOne(filter, &a); err != nil {
			return internal.Author{}, errors.Wrapf(err, "db - author not found")
		}
		return a, nil
	}
}
