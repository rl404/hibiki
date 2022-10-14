package cron

import (
	"context"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rl404/hibiki/internal/errors"
	"github.com/rl404/hibiki/internal/utils"
)

// Fill to fill missing manga.
func (c *Cron) Fill(nrApp *newrelic.Application, limit int) error {
	ctx := errors.Init(context.Background())
	defer c.log(ctx)

	tx := nrApp.StartTransaction("Cron fill")
	defer tx.End()

	ctx = newrelic.NewContext(ctx, tx)

	if err := c.queueMissingManga(ctx, nrApp, limit); err != nil {
		tx.NoticeError(err)
		return errors.Wrap(ctx, err)
	}

	return nil
}

func (c *Cron) queueMissingManga(ctx context.Context, nrApp *newrelic.Application, limit int) error {
	defer newrelic.FromContext(ctx).StartSegment("queueMissingManga").End()

	cnt, _, err := c.service.QueueMissingManga(ctx, limit)
	if err != nil {
		return errors.Wrap(ctx, err)
	}

	utils.Info("queued %d manga", cnt)
	nrApp.RecordCustomEvent("QueueMissingManga", map[string]interface{}{"count": cnt})

	return nil
}
