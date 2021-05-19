package internal

import (
	"github.com/mmuoDev/commons/uuid"
)

//Author is the internal representation of an author
type Author struct {
	ID        uuid.V4 `bson:"id"`
	Pseudonym string  `bson:"pseudonym"`
	Username  string  `bson:"username"`
	Password  string  `bson:"password"`
}

//Book is the internal representation of a book
type Book struct {
	ID          uuid.V4 `bson:"id"`
	AuthorID    string  `bson:"author_id"`
	Title       string  `bson:"title"`
	Description string  `bson:"description"`
	CoverImage  string  `bson:"coverImage"`
	Price       float64 `bson:"price"`
}

//Token represents a token
type Token struct {
	Access string
}
