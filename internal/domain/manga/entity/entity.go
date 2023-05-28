package entity

import "time"

// Manga is entity for manga.
type Manga struct {
	ID                int64
	Title             string
	AlternativeTitles AlternativeTitle
	Picture           string
	StartDate         Date
	EndDate           Date
	Synopsis          string
	Background        string
	NSFW              bool
	Type              Type
	Status            Status
	Chapter           int
	Volume            int
	Mean              float64
	Rank              int
	Popularity        int
	Member            int
	Voter             int
	Favorite          int
	Genres            []Genre
	Pictures          []string
	Related           []Related
	Authors           []Author
	Serialization     []Magazine
	UpdatedAt         time.Time
}

// AlternativeTitle is entity for alternative title.
type AlternativeTitle struct {
	Synonyms []string
	English  string
	Japanese string
}

// Date is entity for date.
type Date struct {
	Day   int
	Month int
	Year  int
}

// Genre is entity for genre.
type Genre struct {
	ID   int64
	Name string
}

// Related is entity for related.
type Related struct {
	ID       int64
	Title    string
	Relation Relation
	Picture  string
}

// Author is entity for author.
type Author struct {
	ID   int64
	Name string
	Role string
}

// Magazine is entity for magazine.
type Magazine struct {
	ID   int64
	Name string
}

// GetAllRequest is get all request model.
type GetAllRequest struct {
	Mode       SearchMode
	Title      string
	Type       Type
	StartDate  *time.Time
	EndDate    *time.Time
	AuthorID   int64
	MagazineID int64
	GenreID    int64
	NSFW       *bool
	Sort       string
	Page       int
	Limit      int
}
