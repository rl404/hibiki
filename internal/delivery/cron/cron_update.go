package cron

import (
	"context"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/hibiki/internal/utils"
)

// Update to update old data.
func (c *Cron) Update(limit int) error {
	ctx := stack.Init(context.Background())
	defer c.log(ctx)

	tx := c.nrApp.StartTransaction("Cron update")
	defer tx.End()

	ctx = newrelic.NewContext(ctx, tx)

	if err := c.queueOldReleasingManga(ctx, limit); err != nil {
		return stack.Wrap(ctx, err)
	}

	if err := c.queueOldFinishedManga(ctx, limit); err != nil {
		return stack.Wrap(ctx, err)
	}

	if err := c.queueOldNotYetManga(ctx, limit); err != nil {
		return stack.Wrap(ctx, err)
	}

	if err := c.queueOldUsername(ctx, limit); err != nil {
		return stack.Wrap(ctx, err)
	}

	return nil
}

func (c *Cron) queueOldReleasingManga(ctx context.Context, limit int) error {
	defer newrelic.FromContext(ctx).StartSegment("queueOldReleasingManga").End()

	cnt, _, err := c.service.QueueOldReleasingManga(ctx, limit)
	if err != nil {
		return stack.Wrap(ctx, err)
	}

	utils.Info("queued %d old releasing manga", cnt)
	c.nrApp.RecordCustomEvent("QueueOldReleasingManga", map[string]interface{}{"count": cnt})

	return nil
}

func (c *Cron) queueOldFinishedManga(ctx context.Context, limit int) error {
	defer newrelic.FromContext(ctx).StartSegment("queueOldFinishedManga").End()

	cnt, _, err := c.service.QueueOldFinishedManga(ctx, limit)
	if err != nil {
		return stack.Wrap(ctx, err)
	}

	utils.Info("queued %d old finished manga", cnt)
	c.nrApp.RecordCustomEvent("QueueOldFinishedManga", map[string]interface{}{"count": cnt})

	return nil
}

func (c *Cron) queueOldNotYetManga(ctx context.Context, limit int) error {
	defer newrelic.FromContext(ctx).StartSegment("queueOldNotYetManga").End()

	cnt, _, err := c.service.QueueOldNotYetManga(ctx, limit)
	if err != nil {
		return stack.Wrap(ctx, err)
	}

	utils.Info("queued %d old not yet released manga", cnt)
	c.nrApp.RecordCustomEvent("QueueOldNotYetManga", map[string]interface{}{"count": cnt})

	return nil
}

func (c *Cron) queueOldUsername(ctx context.Context, limit int) error {
	defer newrelic.FromContext(ctx).StartSegment("queueOldUsername").End()

	cnt, _, err := c.service.QueueOldUserManga(ctx, limit)
	if err != nil {
		return stack.Wrap(ctx, err)
	}

	utils.Info("queued %d old username", cnt)
	c.nrApp.RecordCustomEvent("QueueOldUsername", map[string]interface{}{"count": cnt})

	return nil
}
