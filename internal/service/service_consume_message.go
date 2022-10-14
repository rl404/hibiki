package service

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/rl404/hibiki/internal/domain/publisher/entity"
	"github.com/rl404/hibiki/internal/errors"
)

// ConsumeMessage to consume message from queue.
// Each message type will be handled differently.
func (s *service) ConsumeMessage(ctx context.Context, data entity.Message) error {
	switch data.Type {
	case entity.TypeParseManga:
		return errors.Wrap(ctx, s.consumeParseManga(ctx, data.Data))
	case entity.TypeParseUserManga:
		return errors.Wrap(ctx, s.consumeParseUserManga(ctx, data.Data))
	default:
		return errors.Wrap(ctx, errors.ErrInvalidMessageType)
	}
}

func (s *service) consumeParseManga(ctx context.Context, data []byte) error {
	var req entity.ParseMangaRequest
	if err := json.Unmarshal(data, &req); err != nil {
		return errors.Wrap(ctx, errors.ErrInvalidRequestFormat)
	}

	if !req.Forced {
		if code, err := s.validateID(ctx, req.ID); err != nil {
			if code == http.StatusNotFound {
				return nil
			}
			return errors.Wrap(ctx, err)
		}

		isOld, _, err := s.manga.IsOld(ctx, req.ID)
		if err != nil {
			return errors.Wrap(ctx, err)
		}

		if !isOld {
			return nil
		}
	} else {
		// Delete existing empty id.
		if _, err := s.emptyID.Delete(ctx, req.ID); err != nil {
			return errors.Wrap(ctx, err)
		}
	}

	if _, err := s.updateManga(ctx, req.ID); err != nil {
		return errors.Wrap(ctx, err)
	}

	return nil
}

func (s *service) consumeParseUserManga(ctx context.Context, data []byte) error {
	var req entity.ParseUserMangaRequest
	if err := json.Unmarshal(data, &req); err != nil {
		return errors.Wrap(ctx, errors.ErrInvalidRequestFormat)
	}

	if !req.Forced {
		isOld, _, err := s.userManga.IsOld(ctx, req.Username)
		if err != nil {
			return errors.Wrap(ctx, err)
		}
		if !isOld {
			return nil
		}
	}

	if _, err := s.updateUserManga(ctx, req.Username); err != nil {
		return errors.Wrap(ctx, err)
	}

	return nil
}
