package utils

import (
	"context"
	"golang.org/x/sync/errgroup"
)

type ErrGroupCountdown struct {
	group *errgroup.Group
	wg    ChannelWaitGroup
}

func NewErrGroupCountdown(context context.Context) *ErrGroupCountdown {
	group, _ := errgroup.WithContext(context)

	return &ErrGroupCountdown{
		group: group,
		wg:    NewChannelWaitGroup(),
	}
}

func (g *ErrGroupCountdown) Go(f func() error) {
	g.wg.Add(1)

	g.group.Go(func() error {
		defer g.wg.Done()
		return f()
	})
}

func (g *ErrGroupCountdown) Countdown() chan struct{} {
	return g.wg.Wait()
}

func (g *ErrGroupCountdown) Wait() error {
	return g.group.Wait()
}