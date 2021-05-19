package app

import (
	"books-api/internal/db"
	"books-api/internal/workflow"
	"books-api/pkg"
	"net/http"

	"github.com/mmuoDev/commons/httputils"
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