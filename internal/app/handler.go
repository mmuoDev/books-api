package app

import (
	"books-api/internal"
	"books-api/internal/db"
	"books-api/internal/workflow"
	"books-api/pkg"
	"net/http"
	"reflect"

	"github.com/gorilla/schema"
	"github.com/julienschmidt/httprouter"
	"github.com/mmuoDev/commons/httputils"
	"github.com/mmuoDev/commons/time"
	"github.com/pkg/errors"
)

//AddAuthorHandler returns a http request to add an author
func AddAuthorHandler(addUser db.AddAuthorFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var author pkg.AuthorRequest
		httputils.JSONToDTO(&author, w, r)

		add := workflow.AddAuthor(addUser)
		if err := add(author); err != nil {
			httputils.ServeError(err, w)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

//AuthenticateHandler returns a http request to authenticate a user
func AuthenticateHandler(retrieveAuthor db.RetrieveAuthorByUsernameFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var a pkg.AuthRequest

		httputils.JSONToDTO(&a, w, r)

		auth := workflow.Authenticate(retrieveAuthor)
		u, err := auth(a)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		httputils.ServeJSON(u, w)
	}
}

//AddBookHandler returns a http request to add a book
func AddBookHandler(addBook db.AddBookFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var br pkg.BookRequest
		httputils.JSONToDTO(&br, w, r)
		//get author
		token, err := internal.GetTokenMetaData(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		aID := token.UserID
		add := workflow.AddBook(addBook)
		if err := add(br, aID); err != nil {
			httputils.ServeError(err, w)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

//RetrieveBooksHandler returns a http request to retrieve books
func RetrieveBooksHandler(retrieveBooks db.RetrieveBooksFunc, retrieveAuthor db.RetrieveAuthorByIDFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var params pkg.QueryParams
		if err := GetQueryParams(&params, r); err != nil {
			httputils.ServeError(err, w)
			return
		}
		retrieve := workflow.RetrieveBooks(retrieveBooks, retrieveAuthor)
		books, err := retrieve(params)
		if err != nil {
			httputils.ServeError(err, w)
			return
		}
		httputils.ServeJSON(books, w)
	}
}

//DeleteBookByIDHandler deletes a book by id
func DeleteBookByIDHandler(deleteBook db.DeleteBookByIDFunc, retrieveBook db.RetrieveBookByAuthorIDFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := internal.GetTokenMetaData(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		aID := token.UserID
		params := httprouter.ParamsFromContext(r.Context())
		bID := params.ByName(bookID)
		//
		b, _ := retrieveBook(token.UserID, bID)
		if (internal.Book{}) == b {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		delete := workflow.DeleteBookByID(deleteBook)
		if err := delete(aID, bID); err != nil {
			httputils.ServeError(err, w)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

//RetrieveBookByIDHandler retriees a book by id
func RetrieveBookByIDHandler(retrieveBook db.RetrieveBookByIDFunc, retrieveAuthor db.RetrieveAuthorByIDFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		retrieve := workflow.RetrieveBookByID(retrieveBook, retrieveAuthor)
		params := httprouter.ParamsFromContext(r.Context())
		bID := params.ByName(bookID)
		book, err := retrieve(bID)
		if err != nil {
			httputils.ServeError(err, w)
			return
		}
		httputils.ServeJSON(book, w)
	}
}

//UpdateBookHandler updates a book by its id
func UpdateBookHandler(updateBook db.UpdateBookFunc, retrieveBook db.RetrieveBookByAuthorIDFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var br pkg.BookUpdateRequest
		httputils.JSONToDTO(&br, w, r)

		token, err := internal.GetTokenMetaData(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		params := httprouter.ParamsFromContext(r.Context())
		bID := params.ByName(bookID)
		//

		b, _ := retrieveBook(token.UserID, bID)
		if (internal.Book{}) == b {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		up := workflow.UpdateBook(updateBook)
		if err := up(bID, br); err != nil {
			httputils.ServeError(err, w)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

var decoder = schema.NewDecoder()

// GetQueryParams maps the query params from an http request into an interface
func GetQueryParams(value interface{}, r *http.Request) error {
	// decoder lookup for values on the json tag, instead of the default schema tag
	decoder.SetAliasTag("json")

	var globalErr error

	// Decoder Register for custom type ISO8601
	decoder.RegisterConverter(time.ISO8601{}, func(input string) reflect.Value {
		ISOTime, errISO := time.NewISO8601(input)

		if errISO != nil {
			globalErr = errors.Wrapf(errISO, "handler - invalid iso time provided")
			return reflect.ValueOf(time.ISO8601{})
		}

		return reflect.ValueOf(ISOTime)
	})

	// Decoder Register for custom type Epoch
	decoder.RegisterConverter(time.Epoch(0), func(input string) reflect.Value {
		ISOTime, errISO := time.NewISO8601(input)

		if errISO != nil {
			globalErr = errors.Wrapf(errISO, "handler - invalid iso time provided")
			return reflect.ValueOf(time.ISO8601{}.ToEpoch())
		}

		return reflect.ValueOf(ISOTime.ToEpoch())
	})

	if err := decoder.Decode(value, r.URL.Query()); err != nil {
		return errors.Wrapf(err, "handler - failed to decode query params")
	}

	if globalErr != nil {
		return globalErr
	}

	return nil
}
