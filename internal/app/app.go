package app

import (
	"books-api/internal/db"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mmuoDev/commons/mongo"
)

const (
	bookID = "book_id"
)

//App contains handlers for the app
type App struct {
	AddAuthorHandler        http.HandlerFunc
	AuthenticateHandler     http.HandlerFunc
	AddBookHandler          http.HandlerFunc
	RetrieveBooksHandler    http.HandlerFunc
	DeleteBookByIDHandler   http.HandlerFunc
	RetrieveBookByIDHandler http.HandlerFunc
	UpdateBookHandler       http.HandlerFunc
}

//Handler returns the main handler for this application
func (a App) Handler() http.HandlerFunc {
	router := httprouter.New()

	//books
	router.HandlerFunc(http.MethodPost, "/books", a.AddBookHandler)
	router.HandlerFunc(http.MethodGet, "/books", a.RetrieveBooksHandler)
	router.HandlerFunc(http.MethodDelete, fmt.Sprintf("/books/:%s", bookID), a.DeleteBookByIDHandler)
	router.HandlerFunc(http.MethodPut, fmt.Sprintf("/books/:%s", bookID), a.UpdateBookHandler)
	router.HandlerFunc(http.MethodGet, fmt.Sprintf("/books/:%s", bookID), a.RetrieveBookByIDHandler)

	//authors
	router.HandlerFunc(http.MethodPost, "/authors", a.AddAuthorHandler)

	//auth
	router.HandlerFunc(http.MethodPost, "/auth", a.AuthenticateHandler)

	return http.HandlerFunc(router.ServeHTTP)
}

// Options is a type for application options to modify the app
type Options func(o *OptionalArgs)

// /OptionalArgs optional arguments for this application
type OptionalArgs struct {
	AddAuthor                db.AddAuthorFunc
	RetrieveAuthorByUsername db.RetrieveAuthorByUsernameFunc
	AddBook                  db.AddBookFunc
	RetrieveBooks            db.RetrieveBooksFunc
	DeleteBookByID           db.DeleteBookByIDFunc
	RetrieveAuthor           db.RetrieveAuthorByIDFunc
	RetrieveBookByID         db.RetrieveBookByIDFunc
	UpdateBook               db.UpdateBookFunc
	RetrieveBookByAuthor     db.RetrieveBookByAuthorIDFunc
}

//New creates a new instance of the App
func New(dbProvider mongo.DbProviderFunc, options ...Options) App {
	o := OptionalArgs{
		AddAuthor:                db.AddAuthor(dbProvider),
		RetrieveAuthorByUsername: db.RetrieveAuthorByUsername(dbProvider),
		AddBook:                  db.AddBook(dbProvider),
		RetrieveBooks:            db.RetrieveBooks(dbProvider),
		DeleteBookByID:           db.DeleteBookByID(dbProvider),
		RetrieveAuthor:           db.RetrieveAuthorByID(dbProvider),
		RetrieveBookByID:         db.RetrieveBookByID(dbProvider),
		UpdateBook:               db.UpdateBook(dbProvider),
		RetrieveBookByAuthor:     db.RetrieveBookyAuthorID(dbProvider),
	}

	for _, option := range options {
		option(&o)
	}

	addAuthor := AddAuthorHandler(o.AddAuthor)
	authenticate := AuthenticateHandler(o.RetrieveAuthorByUsername)
	addBook := AddBookHandler(o.AddBook)
	retrieveBooks := RetrieveBooksHandler(o.RetrieveBooks, o.RetrieveAuthor)
	deleteBookByID := DeleteBookByIDHandler(o.DeleteBookByID, o.RetrieveBookByAuthor)
	retrieveBookByID := RetrieveBookByIDHandler(o.RetrieveBookByID, o.RetrieveAuthor)
	updateBook := UpdateBookHandler(o.UpdateBook, o.RetrieveBookByAuthor)

	return App{
		AddAuthorHandler:        addAuthor,
		AuthenticateHandler:     authenticate,
		AddBookHandler:          addBook,
		RetrieveBooksHandler:    retrieveBooks,
		DeleteBookByIDHandler:   deleteBookByID,
		RetrieveBookByIDHandler: retrieveBookByID,
		UpdateBookHandler:       updateBook,
	}
}
