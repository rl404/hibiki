package entity

type messageType string

// Available message type.
const (
	TypeParseManga     messageType = "parse-manga"
	TypeParseUserManga messageType = "parse-user-manga"
)

// Message is entity for message.
type Message struct {
	Type messageType `json:"type"`
	Data []byte      `json:"data"`
}

// ParseMangaRequest is parse manga request model.
type ParseMangaRequest struct {
	ID     int64 `json:"id"`
	Forced bool  `json:"forced"`
}

// ParseUserMangaRequest is parse user manga request model.
type ParseUserMangaRequest struct {
	Username string `json:"username"`
	Forced   bool   `json:"forced"`
}
