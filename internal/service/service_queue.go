package service

import (
	"context"
	"net/http"

	"github.com/rl404/fairy/errors/stack"
)

// QueueOldReleasingManga to queue old releasing manga data.
func (s *service) QueueOldReleasingManga(ctx context.Context, limit int) (int, int, error) {
	var cnt int

	ids, code, err := s.manga.GetOldReleasingIDs(ctx)
	if err != nil {
		return cnt, code, stack.Wrap(ctx, err)
	}

	for i := 0; i < len(ids) && cnt < limit; i, cnt = i+1, cnt+1 {
		if err := s.publisher.PublishParseManga(ctx, ids[i]); err != nil {
			return cnt, http.StatusInternalServerError, stack.Wrap(ctx, err)
		}
	}

	return cnt, http.StatusOK, nil
}

// QueueOldFinishedManga to queue old finished manga data.
func (s *service) QueueOldFinishedManga(ctx context.Context, limit int) (int, int, error) {
	var cnt int

	ids, code, err := s.manga.GetOldFinishedIDs(ctx)
	if err != nil {
		return cnt, code, stack.Wrap(ctx, err)
	}

	for i := 0; i < len(ids) && cnt < limit; i, cnt = i+1, cnt+1 {
		if err := s.publisher.PublishParseManga(ctx, ids[i]); err != nil {
			return cnt, http.StatusInternalServerError, stack.Wrap(ctx, err)
		}
	}

	return cnt, http.StatusOK, nil
}

// QueueOldNotYetManga to queue old not yet released manga data.
func (s *service) QueueOldNotYetManga(ctx context.Context, limit int) (int, int, error) {
	var cnt int

	ids, code, err := s.manga.GetOldFinishedIDs(ctx)
	if err != nil {
		return cnt, code, stack.Wrap(ctx, err)
	}

	for i := 0; i < len(ids) && cnt < limit; i, cnt = i+1, cnt+1 {
		if err := s.publisher.PublishParseManga(ctx, ids[i]); err != nil {
			return cnt, http.StatusInternalServerError, stack.Wrap(ctx, err)
		}
	}

	return cnt, http.StatusOK, nil
}

// QueueMissingManga to queue missing manga.
func (s *service) QueueMissingManga(ctx context.Context, limit int) (int, int, error) {
	var cnt int

	// Get max id.
	maxID, code, err := s.manga.GetMaxID(ctx)
	if err != nil {
		return cnt, code, stack.Wrap(ctx, err)
	}

	// Get all existing manga id.
	mangaIDs, code, err := s.manga.GetIDs(ctx)
	if err != nil {
		return cnt, code, stack.Wrap(ctx, err)
	}

	// Get all empty manga id,
	emptyIDs, code, err := s.emptyID.GetIDs(ctx)
	if err != nil {
		return cnt, code, stack.Wrap(ctx, err)
	}

	idMap := make(map[int64]bool)
	for _, id := range mangaIDs {
		idMap[id] = true
	}
	for _, id := range emptyIDs {
		idMap[id] = true
	}

	// Loop until max id.
	for id := int64(1); id <= maxID && cnt < limit; id++ {
		if idMap[id] {
			continue
		}

		if err := s.publisher.PublishParseManga(ctx, id); err != nil {
			return cnt, http.StatusInternalServerError, stack.Wrap(ctx, err)
		}

		cnt++
	}

	return cnt, http.StatusOK, nil
}

// QueueOldUserManga to queue old user manga.
func (s *service) QueueOldUserManga(ctx context.Context, limit int) (int, int, error) {
	var cnt int

	usernames, code, err := s.userManga.GetOldUsernames(ctx)
	if err != nil {
		return cnt, code, stack.Wrap(ctx, err)
	}

	for i := 0; i < len(usernames) && cnt < limit; i, cnt = i+1, cnt+1 {
		if err := s.publisher.PublishParseUserManga(ctx, usernames[i]); err != nil {
			return cnt, http.StatusInternalServerError, stack.Wrap(ctx, err)
		}
	}

	return cnt, http.StatusOK, nil
}
