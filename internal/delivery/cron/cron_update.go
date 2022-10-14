package cron

import (
	"context"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rl404/hibiki/internal/errors"
	"github.com/rl404/hibiki/internal/utils"
)

// Update to update old data.
func (c *Cron) Update(nrApp *newrelic.Application, limit int) error {
	ctx := errors.Init(context.Background())
	defer c.log(ctx)

	tx := nrApp.StartTransaction("Cron update")
	defer tx.End()

	ctx = newrelic.NewContext(ctx, tx)

	if err := c.queueOldReleasingManga(ctx, nrApp, limit); err != nil {
		tx.NoticeError(err)
		return errors.Wrap(ctx, err)
	}

	if err := c.queueOldFinishedManga(ctx, nrApp, limit); err != nil {
		tx.NoticeError(err)
		return errors.Wrap(ctx, err)
	}

	if err := c.queueOldNotYetManga(ctx, nrApp, limit); err != nil {
		tx.NoticeError(err)
		return errors.Wrap(ctx, err)
	}

	if err := c.queueOldUsername(ctx, nrApp, limit); err != nil {
		tx.NoticeError(err)
		return errors.Wrap(ctx, err)
	}

	return nil
}

func (c *Cron) queueOldReleasingManga(ctx context.Context, nrApp *newrelic.Application, limit int) error {
	defer newrelic.FromContext(ctx).StartSegment("queueOldReleasingManga").End()

	cnt, _, err := c.service.QueueOldReleasingManga(ctx, limit)
	if err != nil {
		return errors.Wrap(ctx, err)
	}

	utils.Info("queued %d old releasing manga", cnt)
	nrApp.RecordCustomEvent("QueueOldReleasingManga", map[string]interface{}{"count": cnt})

	return nil
}

func (c *Cron) queueOldFinishedManga(ctx context.Context, nrApp *newrelic.Application, limit int) error {
	defer newrelic.FromContext(ctx).StartSegment("queueOldFinishedManga").End()

	cnt, _, err := c.service.QueueOldFinishedManga(ctx, limit)
	if err != nil {
		return errors.Wrap(ctx, err)
	}

	utils.Info("queued %d old finished manga", cnt)
	nrApp.RecordCustomEvent("QueueOldFinishedManga", map[string]interface{}{"count": cnt})

	return nil
}

func (c *Cron) queueOldNotYetManga(ctx context.Context, nrApp *newrelic.Application, limit int) error {
	defer newrelic.FromContext(ctx).StartSegment("queueOldNotYetManga").End()

	cnt, _, err := c.service.QueueOldNotYetManga(ctx, limit)
	if err != nil {
		return errors.Wrap(ctx, err)
	}

	utils.Info("queued %d old not yet released manga", cnt)
	nrApp.RecordCustomEvent("QueueOldNotYetManga", map[string]interface{}{"count": cnt})

	return nil
}

func (c *Cron) queueOldUsername(ctx context.Context, nrApp *newrelic.Application, limit int) error {
	defer newrelic.FromContext(ctx).StartSegment("queueOldUsername").End()

	cnt, _, err := c.service.QueueOldUserManga(ctx, limit)
	if err != nil {
		return errors.Wrap(ctx, err)
	}

	utils.Info("queued %d old username", cnt)
	nrApp.RecordCustomEvent("QueueOldUsername", map[string]interface{}{"count": cnt})

	return nil
}
