package entity

import "time"

// UserManga is entity for user manga.
type UserManga struct {
	Username    string
	MangaID     int64
	Title       string
	Status      Status
	Score       int
	Volume      int
	Chapter     int
	StartDate   Date
	EndDate     Date
	Priority    Priority
	IsRereading bool
	RereadCount int
	RereadValue RereadValue
	Tags        []string
	Comment     string
	UpdatedAt   time.Time
}

// Date is entity for date.
type Date struct {
	Year  int
	Month int
	Day   int
}

// GetUserMangaRequest is get user manga request model.
type GetUserMangaRequest struct {
	Username string
	Page     int
	Limit    int
}
