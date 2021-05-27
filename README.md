# books-api

APIs to add authors, authenticate authors and perform CRUD operations on a books resource. This was written with Golang and MongoDB

## Requirements
- Golang
- MongoDB

## Functionalities
- Add Author
- Authenticate Author
- Add Book
- Retrieve Books (with query params)
- Delete Book
- Update a Book
- Get a Book

## Usage
Clone project and `cd` into project foler

### Starting server
``` bash
$ make run
```  

### Running Tests
``` bash
$ make test
```  

### OpenAPI Spec
The OpenAPI spec for this service can be found in the `open-api-yaml` file. Upload to https://editor.swagger.io/ to view. 

### Postman Collection
The Postman collection can be found in the `postman` folder
