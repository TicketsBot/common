package rpc

import (
	"context"
	"github.com/TicketsBot/common/utils"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/atomic"
	"go.uber.org/zap"
)

type Client struct {
	config Config
	client *kgo.Client
	logger *zap.Logger

	consumerRunning *atomic.Bool
	listeners       map[string]Listener

	cancelFunc context.CancelFunc
}

type Config struct {
	Brokers             []string
	ConsumerGroup       string
	ConsumerConcurrency int
}

func NewClient(logger *zap.Logger, config Config, listeners map[string]Listener) (*Client, error) {
	kafkaClient, err := connectKafka(config.Brokers, config.ConsumerGroup, utils.Keys(listeners))
	if err != nil {
		return nil, err
	}

	return &Client{
		config:          config,
		client:          kafkaClient,
		logger:          logger,
		consumerRunning: atomic.NewBool(false),
		listeners:       listeners,
	}, nil
}

func (c *Client) Shutdown() {
	c.client.Close()

	if c.cancelFunc != nil {
		c.cancelFunc()
	}
}

func connectKafka(brokers []string, consumerGroup string, topics []string) (*kgo.Client, error) {
	return kgo.NewClient(
		kgo.SeedBrokers(brokers...),
		kgo.ConsumerGroup(consumerGroup),
		kgo.ConsumeTopics(topics...),
		kgo.ConsumeResetOffset(kgo.NewOffset().AtEnd()),
	)
}
