package cron

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/hibiki/internal/service"
	"github.com/rl404/hibiki/internal/utils"
	"github.com/rl404/hibiki/pkg/log"
)

// Cron contains functions for cron.
type Cron struct {
	service service.Service
	nrApp   *newrelic.Application
}

// New to create new cron.
func New(service service.Service, nrApp *newrelic.Application) *Cron {
	return &Cron{
		service: service,
		nrApp:   nrApp,
	}
}

func (c *Cron) log(ctx context.Context) {
	if rvr := recover(); rvr != nil {
		_ = stack.Wrap(ctx, fmt.Errorf("%s", debug.Stack()), fmt.Errorf("%v", rvr), fmt.Errorf("panic"))
	}

	errStack := stack.Get(ctx)
	if len(errStack) > 0 {
		utils.Log(map[string]interface{}{
			"level": log.ErrorLevel,
			"error": errStack,
		})
	}
}
