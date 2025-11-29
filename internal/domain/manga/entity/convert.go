package entity

import (
	"context"

	"github.com/rl404/nagato"
)

// MangaFromNagato to convert nagato to manga.
func MangaFromNagato(ctx context.Context, manga *nagato.Manga) (*Manga, error) {
	picture := manga.MainPicture.Large
	if picture == "" {
		picture = manga.MainPicture.Medium
	}

	genres := make([]Genre, len(manga.Genres))
	for i, g := range manga.Genres {
		genres[i] = Genre{
			ID:   int64(g.ID),
			Name: g.Name,
		}
	}

	pictures := make([]string, len(manga.Pictures))
	for i, p := range manga.Pictures {
		pictures[i] = p.Large
		if pictures[i] == "" {
			pictures[i] = p.Medium
		}
	}

	related := make([]Related, len(manga.RelatedManga))
	for i, r := range manga.RelatedManga {
		pic := r.Manga.MainPicture.Large
		if pic == "" {
			pic = r.Manga.MainPicture.Medium
		}
		related[i] = Related{
			ID:       int64(r.Manga.ID),
			Title:    r.Manga.Title,
			Relation: nagatoToRelation(r.RelationType),
			Picture:  pic,
		}
	}

	authors := make([]Author, len(manga.Authors))
	for i, a := range manga.Authors {
		authors[i] = Author{
			ID:   int64(a.Person.ID),
			Name: a.Person.FirstName + " " + a.Person.LastName,
			Role: a.Role,
		}
	}

	magazines := make([]Magazine, len(manga.Serialization))
	for i, m := range manga.Serialization {
		magazines[i] = Magazine{
			ID:   int64(m.Magazine.ID),
			Name: m.Magazine.Name,
		}
	}

	return &Manga{
		ID:    int64(manga.ID),
		Title: manga.Title,
		AlternativeTitles: AlternativeTitle{
			Synonyms: manga.AlternativeTitles.Synonyms,
			English:  manga.AlternativeTitles.English,
			Japanese: manga.AlternativeTitles.Japanese,
		},
		Picture: picture,
		StartDate: Date{
			Year:  manga.StartDate.Year,
			Month: manga.StartDate.Month,
			Day:   manga.StartDate.Day,
		},
		EndDate: Date{
			Year:  manga.EndDate.Year,
			Month: manga.EndDate.Month,
			Day:   manga.EndDate.Day,
		},
		Synopsis:      manga.Synopsis,
		Background:    manga.Background,
		NSFW:          manga.NSFW != nagato.NsfwWhite,
		Type:          nagatoToType(manga.MediaType),
		Status:        nagatoToStatus(manga.Status),
		Chapter:       manga.NumChapters,
		Volume:        manga.NumVolumes,
		Mean:          manga.Mean,
		Rank:          manga.Rank,
		Popularity:    manga.Popularity,
		Member:        manga.NumListUsers,
		Voter:         manga.NumScoringUsers,
		Favorite:      manga.NumFavorites,
		Genres:        genres,
		Pictures:      pictures,
		Related:       related,
		Authors:       authors,
		Serialization: magazines,
	}, nil
}

func nagatoToType(t nagato.MediaType) Type {
	return map[nagato.MediaType]Type{
		nagato.MediaManga:      TypeManga,
		nagato.MediaNovel:      TypeNovel,
		nagato.MediaOneShot:    TypeOneShot,
		nagato.MediaDoujinshi:  TypeDoujinshi,
		nagato.MediaManhwa:     TypeManhwa,
		nagato.MediaManhua:     TypeManhua,
		nagato.MediaOEL:        TypeOEL,
		nagato.MediaLightNovel: TypeLightNovel,
	}[t]
}

func nagatoToStatus(t nagato.StatusType) Status {
	return map[nagato.StatusType]Status{
		nagato.StatusFinishedPublishing:  StatusFinished,
		nagato.StatusCurrentlyPublishing: StatusReleasing,
		nagato.StatusNotYetPublished:     StatusNotYet,
		nagato.StatusOnHiatus:            StatusHiatus,
		nagato.StatusDiscontinued:        StatusDiscontinued,
	}[t]
}

func nagatoToRelation(t nagato.RelationType) Relation {
	return map[nagato.RelationType]Relation{
		nagato.RelationSequel:             RelationSequel,
		nagato.RelationPrequel:            RelationPrequel,
		nagato.RelationAlternativeSetting: RelationAlternativeSetting,
		nagato.RelationAlternativeVersion: RelationAlternativeVersion,
		nagato.RelationSideStory:          RelationSideStory,
		nagato.RelationParentStory:        RelationParentStory,
		nagato.RelationSummary:            RelationSummary,
		nagato.RelationFullStory:          RelationFullStory,
		nagato.RelationSpinOff:            RelationSpinOff,
		nagato.RelationOther:              RelationOther,
		nagato.RelationCharacter:          RelationCharacter,
		nagato.RelationAdaptation:         RelationAdaptation,
	}[t]
}
