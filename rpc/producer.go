package rpc

import (
	"context"
	"encoding/json"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"
)

func (c *Client) ProduceSync(ctx context.Context, topic string, message []byte) error {
	c.logger.Debug("Producing message", zap.String("topic", topic), zap.ByteString("message", message))

	return c.client.ProduceSync(ctx, &kgo.Record{
		Topic: topic,
		Value: message,
	}).FirstErr()
}

func (c *Client) ProduceSyncJson(ctx context.Context, topic string, message any) error {
	bytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return c.ProduceSync(ctx, topic, bytes)
}
