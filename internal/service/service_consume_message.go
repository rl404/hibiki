package service

import (
	"context"

	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/hibiki/internal/domain/publisher/entity"
	"github.com/rl404/hibiki/internal/errors"
)

// ConsumeMessage to consume message from queue.
// Each message type will be handled differently.
func (s *service) ConsumeMessage(ctx context.Context, data entity.Message) error {
	switch data.Type {
	case entity.TypeParseManga:
		return stack.Wrap(ctx, s.consumeParseManga(ctx, data))
	case entity.TypeParseUserManga:
		return stack.Wrap(ctx, s.consumeParseUserManga(ctx, data))
	default:
		return stack.Wrap(ctx, errors.ErrInvalidMessageType)
	}
}

func (s *service) consumeParseManga(ctx context.Context, data entity.Message) error {
	if !data.Forced {
		isOld, _, err := s.manga.IsOld(ctx, data.ID)
		if err != nil {
			return stack.Wrap(ctx, err)
		}

		if !isOld {
			return nil
		}
	}

	// Delete existing empty id.
	if _, err := s.emptyID.Delete(ctx, data.ID); err != nil {
		return stack.Wrap(ctx, err)
	}

	if _, err := s.updateManga(ctx, data.ID); err != nil {
		return stack.Wrap(ctx, err)
	}

	return nil
}

func (s *service) consumeParseUserManga(ctx context.Context, data entity.Message) error {
	if !data.Forced {
		isOld, _, err := s.userManga.IsOld(ctx, data.Username)
		if err != nil {
			return stack.Wrap(ctx, err)
		}

		if !isOld {
			return nil
		}
	}

	if _, err := s.updateUserManga(ctx, data.Username); err != nil {
		return stack.Wrap(ctx, err)
	}

	return nil
}
