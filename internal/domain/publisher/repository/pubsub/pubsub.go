package pubsub

import (
	"context"
	"encoding/json"

	"github.com/rl404/fairy/pubsub"
	"github.com/rl404/hibiki/internal/domain/publisher/entity"
	"github.com/rl404/hibiki/internal/errors"
)

// Pubsub contains functions for pubsub.
type Pubsub struct {
	pubsub pubsub.PubSub
	topic  string
}

// New to create new pubsub.
func New(ps pubsub.PubSub, topic string) *Pubsub {
	return &Pubsub{
		pubsub: ps,
		topic:  topic,
	}
}

// PublishParseManga to publish parse manga.
func (p *Pubsub) PublishParseManga(ctx context.Context, data entity.ParseMangaRequest) error {
	d, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(ctx, errors.ErrInternalServer, err)
	}

	if err := p.pubsub.Publish(ctx, p.topic, entity.Message{
		Type: entity.TypeParseManga,
		Data: d,
	}); err != nil {
		return errors.Wrap(ctx, errors.ErrInternalServer, err)
	}

	return nil
}

// PublishParseUserManga to publish parse user manga.
func (p *Pubsub) PublishParseUserManga(ctx context.Context, data entity.ParseUserMangaRequest) error {
	d, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(ctx, errors.ErrInternalServer, err)
	}

	if err := p.pubsub.Publish(ctx, p.topic, entity.Message{
		Type: entity.TypeParseUserManga,
		Data: d,
	}); err != nil {
		return errors.Wrap(ctx, errors.ErrInternalServer, err)
	}

	return nil
}
