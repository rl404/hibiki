package entity

type messageType string

// Available message type.
const (
	TypeParseManga     messageType = "parse-manga"
	TypeParseUserManga messageType = "parse-user-manga"
)

// Message is entity for message.
type Message struct {
	Type     messageType `json:"type"`
	ID       int64       `json:"id"`
	Username string      `json:"username"`
	Forced   bool        `json:"forced"`
}
