# books-api

This is a books API written with Golang

## Requirements
- Golang
- MongoDB

## Functionalities
- Add Author
- Authenticate Author
- Add Book
- Retrieve Books (with query params)
- Delete Book

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
## More details

These endpoints require the bearer token generated during authenticattion
- Add Book
- Delete Book
