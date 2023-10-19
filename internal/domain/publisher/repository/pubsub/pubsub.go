package pubsub

import (
	"context"
	"encoding/json"

	"github.com/rl404/fairy/errors/stack"
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
func (p *Pubsub) PublishParseManga(ctx context.Context, id int64) error {
	d, err := json.Marshal(entity.Message{
		Type: entity.TypeParseManga,
		ID:   id,
	})
	if err != nil {
		return stack.Wrap(ctx, err, errors.ErrInternalServer)
	}

	if err := p.pubsub.Publish(ctx, p.topic, d); err != nil {
		return stack.Wrap(ctx, err, errors.ErrInternalServer)
	}

	return nil
}

// PublishParseUserManga to publish parse user manga.
func (p *Pubsub) PublishParseUserManga(ctx context.Context, username string) error {
	d, err := json.Marshal(entity.Message{
		Type:     entity.TypeParseUserManga,
		Username: username,
	})
	if err != nil {
		return stack.Wrap(ctx, err, errors.ErrInternalServer)
	}

	if err := p.pubsub.Publish(ctx, p.topic, d); err != nil {
		return stack.Wrap(ctx, err, errors.ErrInternalServer)
	}

	return nil
}
