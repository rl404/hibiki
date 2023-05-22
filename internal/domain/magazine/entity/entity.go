package entity

// Magazine is entity for magazine.
type Magazine struct {
	ID   int64
	Name string
}

// GetAllRequest is get all request model.
type GetAllRequest struct {
	Name  string
	Page  int
	Limit int
}
