package rpc

import (
	"context"
	"errors"
	"github.com/panjf2000/ants/v2"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"
)

const maxEventsPerPoll = 100

func (c *Client) StartConsumer() {
	if c.consumerRunning.Swap(true) {
		c.logger.Fatal("Kafka client already running")
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	c.cancelFunc = cancel

	pool, err := ants.NewPool(c.config.ConsumerConcurrency)
	if err != nil {
		c.logger.Fatal("Failed to create worker pool", zap.Error(err))
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			records, err := c.poll(ctx)
			if err != nil {
				if errors.Is(err, kgo.ErrClientClosed) {
					c.logger.Info("Kafka client closed, stopping read loop")
					return
				} else if errors.Is(err, context.Canceled) {
					c.logger.Info("Context cancelled, stopping read loop")
					return
				} else {
					c.logger.Error("Failed to poll records", zap.Error(err))
					continue
				}
			}

			for _, record := range records {
				listener, ok := c.listeners[record.Topic]
				if !ok {
					c.logger.Warn("No listener found for topic", zap.String("topic", record.Topic))
					continue
				}

				value := record.Value
				if err := pool.Submit(func() {
					ctx, cancel := listener.BuildContext()
					defer cancel()

					listener.HandleMessage(ctx, value)
				}); err != nil {
					c.logger.Error("Failed to submit task to worker pool", zap.Error(err))
					continue
				}
			}
		}
	}
}

func (c *Client) poll(ctx context.Context) ([]*kgo.Record, error) {
	fetches := c.client.PollRecords(ctx, maxEventsPerPoll)
	if fetches.IsClientClosed() {
		return nil, kgo.ErrClientClosed
	}

	if err := fetches.Err(); err != nil {
		return nil, err
	}

	records := make([]*kgo.Record, 0, fetches.NumRecords())

	iter := fetches.RecordIter()
	for !iter.Done() {
		record := iter.Next()
		records = append(records, record)
	}

	return records, nil
}
