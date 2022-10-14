package entity

import (
	"github.com/rl404/nagato"
)

// UserMangaFromNagato to convert nagato to user manga.
func UserMangaFromNagato(username string, manga nagato.UserManga) UserManga {
	return UserManga{
		Username: username,
		MangaID:  int64(manga.Manga.ID),
		Title:    manga.Manga.Title,
		Status:   nagatoToStatus(manga.Status.Status),
		Score:    manga.Status.Score,
		Volume:   manga.Status.NumVolumesRead,
		Chapter:  manga.Status.NumChaptersRead,
		StartDate: Date{
			Year:  manga.Status.StartDate.Year,
			Month: manga.Status.StartDate.Month,
			Day:   manga.Status.StartDate.Day,
		},
		EndDate: Date{
			Year:  manga.Status.FinishDate.Year,
			Month: manga.Status.FinishDate.Month,
			Day:   manga.Status.FinishDate.Day,
		},
		Priority:    nagatoToPriority(manga.Status.Priority),
		IsRereading: manga.Status.IsRereading,
		RereadCount: manga.Status.NumTimesReread,
		RereadValue: nagatoToRereadValue(manga.Status.RereadValue),
		Tags:        manga.Status.Tags,
		Comment:     manga.Status.Comments,
	}
}

func nagatoToStatus(t nagato.UserMangaStatusType) Status {
	return map[nagato.UserMangaStatusType]Status{
		nagato.UserMangaStatusReading:    StatusReading,
		nagato.UserMangaStatusCompleted:  StatusCompleted,
		nagato.UserMangaStatusOnHold:     StatusOnHold,
		nagato.UserMangaStatusDropped:    StatusDropped,
		nagato.UserMangaStatusPlanToRead: StatusPlanned,
	}[t]
}

func nagatoToPriority(t nagato.PriorityType) Priority {
	return map[nagato.PriorityType]Priority{
		nagato.PriorityLow:    PriorityLow,
		nagato.PriorityMedium: PriorityMedium,
		nagato.PriorityHigh:   PriorityHigh,
	}[t]
}

func nagatoToRereadValue(t nagato.RereadValueType) RereadValue {
	return map[nagato.RereadValueType]RereadValue{
		nagato.RereadValueVeryLow:  RereadValueVeryLow,
		nagato.RereadValueLow:      RereadValueLow,
		nagato.RereadValueMedium:   RereadValueMedium,
		nagato.RereadValueHigh:     RereadValueHigh,
		nagato.RereadValueVeryHigh: RereadValueVeryHigh,
	}[t]
}
