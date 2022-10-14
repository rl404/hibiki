package entity

import "github.com/rl404/nagato"

// GetUserMangaRequest is get user manga request entity.
type GetUserMangaRequest struct {
	Username string
	Status   nagato.UserMangaStatusType
	Sort     nagato.UserMangaSortType
	Limit    int
	Offset   int
}
