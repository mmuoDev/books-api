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
}

//Handler returns the main handler for this application
func (a App) Handler() http.HandlerFunc {
	router := httprouter.New()

	router.HandlerFunc(http.MethodPost, "/authors", a.AddAuthorHandler)
	router.HandlerFunc(http.MethodPost, "/books", a.AddBookHandler)
	router.HandlerFunc(http.MethodGet, "/books", a.RetrieveBooksHandler)
	router.HandlerFunc(http.MethodDelete, fmt.Sprintf("/books/:%s", bookID), a.DeleteBookByIDHandler)
	router.HandlerFunc(http.MethodGet, fmt.Sprintf("/books/:%s", bookID), a.RetrieveBookByIDHandler)

	router.HandlerFunc(http.MethodPost, "/auth", a.AuthenticateHandler)

	return http.HandlerFunc(router.ServeHTTP)
}

// Options is a type for application options to modify the app
type Options func(o *Option)

// /OptionalArgs optional arguments for this application
type Option struct {
	AddAuthor                db.AddAuthorFunc
	RetrieveAuthorByUsername db.RetrieveAuthorByUsernameFunc
	AddBook                  db.AddBookFunc
	RetrieveBooks            db.RetrieveBooksFunc
	DeleteBookByID           db.DeleteBookByIDFunc
	RetrieveAuthor           db.RetrieveAuthorByIDFunc
	RetrieveBookByID         db.RetrieveBookByIDFunc
}

//New creates a new instance of the App
func New(dbProvider mongo.DbProviderFunc, options ...Options) App {
	o := Option{
		AddAuthor:                db.AddAuthor(dbProvider),
		RetrieveAuthorByUsername: db.RetrieveAuthorByUsername(dbProvider),
		AddBook:                  db.AddBook(dbProvider),
		RetrieveBooks:            db.RetrieveBooks(dbProvider),
		DeleteBookByID:           db.DeleteBookByID(dbProvider),
		RetrieveAuthor:           db.RetrieveAuthorByID(dbProvider),
		RetrieveBookByID:         db.RetrieveBookByID(dbProvider),
	}

	for _, option := range options {
		option(&o)
	}

	addAuthor := AddAuthorHandler(o.AddAuthor)
	authenticate := AuthenticateHandler(o.RetrieveAuthorByUsername)
	addBook := AddBookHandler(o.AddBook)
	retrieveBooks := RetrieveBooksHandler(o.RetrieveBooks, o.RetrieveAuthor)
	deleteBookByID := DeleteBookByIDHandler(o.DeleteBookByID)
	retrieveBookByID := RetrieveBookByIDHandler(o.RetrieveBookByID, o.RetrieveAuthor)

	return App{
		AddAuthorHandler:        addAuthor,
		AuthenticateHandler:     authenticate,
		AddBookHandler:          addBook,
		RetrieveBooksHandler:    retrieveBooks,
		DeleteBookByIDHandler:   deleteBookByID,
		RetrieveBookByIDHandler: retrieveBookByID,
	}
}
