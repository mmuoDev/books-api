package app

import (
	"books-api/internal/db"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mmuoDev/commons/mongo"
)

//App contains handlers for the app
type App struct {
	AddAuthorHandler    http.HandlerFunc
	AuthenticateHandler http.HandlerFunc
}

//Handler returns the main handler for this application
func (a App) Handler() http.HandlerFunc {
	router := httprouter.New()

	router.HandlerFunc(http.MethodPost, "/authors", a.AddAuthorHandler)
	router.HandlerFunc(http.MethodPost, "/auth", a.AuthenticateHandler)

	return http.HandlerFunc(router.ServeHTTP)
}

// Options is a type for application options to modify the app
type Options func(o *Option)

// /OptionalArgs optional arguments for this application
type Option struct {
	AddAuthor                db.AddAuthorFunc
	RetrieveAuthorByUsername db.RetrieveAuthorByUsernameFunc
}

//New creates a new instance of the App
func New(dbProvider mongo.DbProviderFunc, options ...Options) App {
	o := Option{
		AddAuthor:                db.AddAuthor(dbProvider),
		RetrieveAuthorByUsername: db.RetrieveAuthorByUsername(dbProvider),
	}

	for _, option := range options {
		option(&o)
	}

	addAuthor := AddAuthorHandler(o.AddAuthor)
	authenticate := AuthenticateHandler(o.RetrieveAuthorByUsername)

	return App{
		AddAuthorHandler:    addAuthor,
		AuthenticateHandler: authenticate,
	}
}
