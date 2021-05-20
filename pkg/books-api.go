package pkg

//AuthorRequest represents request for creating an author
type AuthorRequest struct {
	Pseudonym string `json:"pseudonym"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

//AuthRequest represents data needed to authenticate a user and generate a token
type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//BookRequest represents request for adding a book for an author
type BookRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	CoverImage  string  `json:"coverImage"`
	Price       float64 `json:"price"`
}

//Book represent details for a book
type Book struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	CoverImage  string  `json:"coverImage"`
	Price       float64 `json:"price"`
	Author      string  `json:"author"`
}

//Auth represents data after successful authentication
type Auth struct {
	AuthorID    string `json:"author_id"`
	AccessToken string `json:"access_token"`
}

// QueryParams represents query params for filtering
type QueryParams struct {
	Title string `json:"title"`
}
