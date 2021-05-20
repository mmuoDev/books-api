package db

import (
	"books-api/internal"
	"books-api/pkg"

	mmuoMongo "github.com/mmuoDev/commons/mongo"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

//AddBookFunc adds a book
type AddBookFunc func(internal.Book) error

//RetrieveBooksFunc retrieves books
type RetrieveBooksFunc func(params pkg.QueryParams) ([]internal.Book, error)

//DeleteBookByIDFunc deletes a book by id
type DeleteBookByIDFunc func(aID, bID string) error

//RetrieveBookByIDFunc retrieves a book by id
type RetrieveBookByIDFunc func(bID string) (internal.Book, error)

//AddAuthor adds an author to DB
func AddAuthor(dbProvider mmuoMongo.DbProviderFunc) AddAuthorFunc {
	return func(a internal.Author) error {
		col := mmuoMongo.NewCollection(dbProvider, authorsCollection)
		_, err := col.Insert(a)
		if err != nil {
			return errors.Wrap(err, "db - failure inserting an author")
		}
		return nil
	}
}

//RetrieveBookByID retrieves book by id
func RetrieveBookByID(dbProvider mmuoMongo.DbProviderFunc) RetrieveBookByIDFunc {
	return func(bID string) (internal.Book, error) {
		col := mmuoMongo.NewCollection(dbProvider, booksCollection)
		var b internal.Book
		if err := col.FindByID(bID, &b); err != nil {
			return b, errors.Wrap(err, "db - failure retrieving book")
		}
		return b, nil
	}
}

//RetrieveAuthorByID retrieves author by the author's id
func RetrieveAuthorByID(dbProvider mmuoMongo.DbProviderFunc) RetrieveAuthorByIDFunc {
	return func(aid string) (internal.Author, error) {
		col := mmuoMongo.NewCollection(dbProvider, authorsCollection)
		var a internal.Author
		if err := col.FindByID(aid, &a); err != nil {
			return a, errors.Wrap(err, "db - failure retrieving author")
		}
		return a, nil
	}
}

//RetrieveAuthorByUsername retrieves author by user
func RetrieveAuthorByUsername(dbProvider mmuoMongo.DbProviderFunc) RetrieveAuthorByUsernameFunc {
	return func(username string) (internal.Author, error) {
		col := mmuoMongo.NewCollection(dbProvider, authorsCollection)
		filter := bson.D{{"username", username}}
		var a internal.Author

		if err := col.FindOne(filter, &a); err != nil {
			return internal.Author{}, errors.Wrapf(err, "db - author not found")
		}
		return a, nil
	}
}

//AddBook adds an author to DB
func AddBook(dbProvider mmuoMongo.DbProviderFunc) AddBookFunc {
	return func(a internal.Book) error {
		col := mmuoMongo.NewCollection(dbProvider, booksCollection)
		_, err := col.Insert(a)
		if err != nil {
			return errors.Wrap(err, "db - failure inserting a book")
		}
		return nil
	}
}

//RetrieveBooks for authenticated user
func RetrieveBooks(dbProvider mmuoMongo.DbProviderFunc) RetrieveBooksFunc {
	return func(params pkg.QueryParams) ([]internal.Book, error) {
		col := mmuoMongo.NewCollection(dbProvider, booksCollection)
		filter := bson.M{}

		if params.Title != "" {
			filter = bson.M{"title": params.Title}
		}

		books := []internal.Book{}
		onEach := func(cur *mongo.Cursor) error {
			b := internal.Book{}
			err := cur.Decode(&b)
			if err != nil {
				return err
			}
			books = append(books, b)
			return nil
		}

		if err := col.FindMulti(filter, onEach); err != nil {
			return nil, errors.Wrapf(err, "db - failure retrieving books")
		}
		return books, nil
	}
}

//DeleteBookByID deletes a book by id
func DeleteBookByID(dbProvider mmuoMongo.DbProviderFunc) DeleteBookByIDFunc {
	return func(aID, bID string) error {
		col := mmuoMongo.NewCollection(dbProvider, booksCollection)
		filter := bson.D{
			{"author_id", aID},
			{"id", bID},
		}
		if err := col.DeleteMany(filter); err != nil {
			errors.Wrapf(err, "db - failure deleting book with id=%s by authorID=%s", bID, aID)
		}
		return nil
	}
}
