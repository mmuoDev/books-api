package app_test

import (
	"books-api/internal"
	"books-api/internal/app"
	"books-api/pkg"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/mmuoDev/commons/httputils"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
)

//mongoDBProvider mocks mongo DB
func mongoDBProvider() *mongo.Database {
	return nil
}

func TestCreateAuthorReturns201(t *testing.T) {
	insertIntoDBInvoked := false

	mockAuthorDBInsert := func(o *app.OptionalArgs) {
		o.AddAuthor = func(a internal.Author) error {
			insertIntoDBInvoked = true
			return nil
		}
	}

	//optional args
	opts := []app.Options{
		mockAuthorDBInsert,
	}

	ap := app.New(mongoDBProvider, opts...)
	serverURL, cleanUpServer := app.NewTestServer(ap.Handler())
	defer cleanUpServer()

	reqPayload, _ := os.Open(filepath.Join("testdata", "add_author_request.json"))
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/authors", serverURL), reqPayload)

	client := &http.Client{}
	res, _ := client.Do(req)

	t.Run("Http Status Code is 201", func(t *testing.T) {
		assert.Equal(t, res.StatusCode, http.StatusCreated)
	})

	t.Run("Insert to DB invoked", func(t *testing.T) {
		assert.True(t, insertIntoDBInvoked)
	})

}

func TestRetrieveBooksReturns200(t *testing.T) {
	retrieveBooksFromDB := false
	retrieveAuthorFromDB := false

	mockRetrieveBooks := func(o *app.OptionalArgs) {
		o.RetrieveBooks = func(params pkg.QueryParams) ([]internal.Book, error) {
			retrieveBooksFromDB = true
			var bRes []internal.Book
			httputils.FileToStruct(filepath.Join("testdata", "retrieve_books_db.json"), &bRes)

			return bRes, nil
		}
	}

	mockRetrieveAuthor := func(o *app.OptionalArgs) {
		o.RetrieveAuthor = func(aid string) (internal.Author, error) {
			retrieveAuthorFromDB = true

			var aRes internal.Author
			httputils.FileToStruct(filepath.Join("testdata", "retrieve_author_db.json"), &aRes)

			return aRes, nil
		}
	}

	//optional args
	opts := []app.Options{
		mockRetrieveAuthor,
		mockRetrieveBooks,
	}

	ap := app.New(mongoDBProvider, opts...)
	serverURL, cleanUpServer := app.NewTestServer(ap.Handler())
	defer cleanUpServer()

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/books", serverURL), nil)

	client := &http.Client{}
	res, _ := client.Do(req)

	t.Run("Http Status Code is 200", func(t *testing.T) {
		assert.Equal(t, res.StatusCode, http.StatusOK)
	})

	t.Run("Retrieve Books from DB invoked", func(t *testing.T) {
		assert.True(t, retrieveBooksFromDB)
	})

	t.Run("Retrieve Author from DB invoked", func(t *testing.T) {
		assert.True(t, retrieveAuthorFromDB)
	})

	t.Run("Response Body is as expected", func(t *testing.T) {
		var (
			expectedResponse []internal.Book
			actualResponse []internal.Book
		)
		json.NewDecoder(res.Body).Decode(&actualResponse)
		httputils.FileToStruct(filepath.Join("testdata", "retrieve_books_response.json"), &expectedResponse)
		assert.Equal(t, expectedResponse, actualResponse)
	})

}
