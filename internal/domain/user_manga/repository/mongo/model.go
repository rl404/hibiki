package mongo

import (
	"time"

	"github.com/rl404/hibiki/internal/domain/user_manga/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type userManga struct {
	Username    string             `bson:"username"`
	MangaID     int64              `bson:"manga_id"`
	Title       string             `bson:"title"`
	Status      entity.Status      `bson:"status"`
	Score       int                `bson:"score"`
	Volume      int                `bson:"volume"`
	Chapter     int                `bson:"chapter"`
	StartDate   date               `bson:"start_date"`
	EndDate     date               `bson:"end_date"`
	Priority    entity.Priority    `bson:"priority"`
	IsRereading bool               `bson:"is_rereading"`
	RereadCount int                `bson:"reread_count"`
	RereadValue entity.RereadValue `bson:"reread_value"`
	Tags        []string           `bson:"tags"`
	Comment     string             `bson:"comment"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
}

func (um *userManga) MarshalBSON() ([]byte, error) {
	if um.CreatedAt.IsZero() {
		um.CreatedAt = time.Now()
	}

	um.UpdatedAt = time.Now()

	type um2 userManga
	return bson.Marshal((*um2)(um))
}

type date struct {
	Year  int `bson:"year"`
	Month int `bson:"month"`
	Day   int `bson:"day"`
}

func (um *userManga) toEntity() *entity.UserManga {
	return &entity.UserManga{
		Username: um.Username,
		MangaID:  um.MangaID,
		Title:    um.Title,
		Status:   um.Status,
		Score:    um.Score,
		Volume:   um.Volume,
		Chapter:  um.Chapter,
		StartDate: entity.Date{
			Year:  um.StartDate.Year,
			Month: um.StartDate.Month,
			Day:   um.StartDate.Day,
		},
		EndDate: entity.Date{
			Year:  um.EndDate.Year,
			Month: um.EndDate.Month,
			Day:   um.EndDate.Day,
		},
		Priority:    um.Priority,
		IsRereading: um.IsRereading,
		RereadCount: um.RereadCount,
		RereadValue: um.RereadValue,
		Tags:        um.Tags,
		Comment:     um.Comment,
		UpdatedAt:   um.UpdatedAt,
	}
}

func (m *Mongo) fromEntity(data entity.UserManga) *userManga {
	return &userManga{
		Username: data.Username,
		MangaID:  data.MangaID,
		Title:    data.Title,
		Status:   data.Status,
		Score:    data.Score,
		Volume:   data.Volume,
		Chapter:  data.Chapter,
		StartDate: date{
			Year:  data.StartDate.Year,
			Month: data.StartDate.Month,
			Day:   data.StartDate.Day,
		},
		EndDate: date{
			Year:  data.EndDate.Year,
			Month: data.EndDate.Month,
			Day:   data.EndDate.Day,
		},
		Priority:    data.Priority,
		IsRereading: data.IsRereading,
		RereadCount: data.RereadCount,
		RereadValue: data.RereadValue,
		Tags:        data.Tags,
		Comment:     data.Comment,
	}
}
