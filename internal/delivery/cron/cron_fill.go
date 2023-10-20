package cron

import (
	"context"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/hibiki/internal/utils"
)

// Fill to fill missing manga.
func (c *Cron) Fill(limit int) error {
	ctx := stack.Init(context.Background())
	defer c.log(ctx)

	tx := c.nrApp.StartTransaction("Cron fill")
	defer tx.End()

	ctx = newrelic.NewContext(ctx, tx)

	if err := c.queueMissingManga(ctx, limit); err != nil {
		return stack.Wrap(ctx, err)
	}

	return nil
}

func (c *Cron) queueMissingManga(ctx context.Context, limit int) error {
	defer newrelic.FromContext(ctx).StartSegment("queueMissingManga").End()

	cnt, _, err := c.service.QueueMissingManga(ctx, limit)
	if err != nil {
		return stack.Wrap(ctx, err)
	}

	utils.Info("queued %d manga", cnt)
	c.nrApp.RecordCustomEvent("QueueMissingManga", map[string]interface{}{"count": cnt})

	return nil
}
