package entity

// Author is entity for author.
type Author struct {
	ID        int64
	FirstName string
	LastName  string
}

// GetAllRequest is get all request model.
type GetAllRequest struct {
	Name  string
	Page  int
	Limit int
}
