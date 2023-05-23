package entity

// Genre is entity for genre.
type Genre struct {
	ID   int64
	Name string
}

// GetAllRequest is get all request model.
type GetAllRequest struct {
	Name  string
	Page  int
	Limit int
}
