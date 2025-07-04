package amqp

import (
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/require"

	"github.com/influxdata/telegraf/config"
)

type MockClient struct {
	PublishF func() error
	CloseF   func() error

	PublishCallCount int
	CloseCallCount   int
}

func (c *MockClient) Publish(string, []byte) error {
	c.PublishCallCount++
	return c.PublishF()
}

func (c *MockClient) Close() error {
	c.CloseCallCount++
	return c.CloseF()
}

func NewMockClient() Client {
	return &MockClient{
		PublishF: func() error {
			return nil
		},
		CloseF: func() error {
			return nil
		},
	}
}

func TestConnect(t *testing.T) {
	tests := []struct {
		name    string
		output  *AMQP
		errFunc func(t *testing.T, output *AMQP, err error)
	}{
		{
			name: "defaults",
			output: &AMQP{
				Brokers:            []string{DefaultURL},
				ExchangeType:       DefaultExchangeType,
				ExchangeDurability: "durable",
				AuthMethod:         DefaultAuthMethod,
				Headers: map[string]string{
					"database":         DefaultDatabase,
					"retention_policy": DefaultRetentionPolicy,
				},
				Timeout: config.Duration(time.Second * 5),
				connect: func(_ *ClientConfig) (Client, error) {
					return NewMockClient(), nil
				},
			},
			errFunc: func(t *testing.T, output *AMQP, err error) {
				cfg := output.config
				require.Equal(t, []string{DefaultURL}, cfg.brokers)
				require.Empty(t, cfg.exchange)
				require.Equal(t, "topic", cfg.exchangeType)
				require.False(t, cfg.exchangePassive)
				require.True(t, cfg.exchangeDurable)
				require.Equal(t, amqp.Table(nil), cfg.exchangeArguments)
				require.Equal(t, amqp.Table{
					"database":         DefaultDatabase,
					"retention_policy": DefaultRetentionPolicy,
				}, cfg.headers)
				require.Equal(t, amqp.Transient, cfg.deliveryMode)
				require.NoError(t, err)
			},
		},
		{
			name: "headers overrides deprecated dbrp",
			output: &AMQP{
				Headers: map[string]string{
					"foo": "bar",
				},
				connect: func(_ *ClientConfig) (Client, error) {
					return NewMockClient(), nil
				},
			},
			errFunc: func(t *testing.T, output *AMQP, err error) {
				cfg := output.config
				require.Equal(t, amqp.Table{
					"foo": "bar",
				}, cfg.headers)
				require.NoError(t, err)
			},
		},
		{
			name: "exchange args",
			output: &AMQP{
				ExchangeArguments: map[string]string{
					"foo": "bar",
				},
				connect: func(_ *ClientConfig) (Client, error) {
					return NewMockClient(), nil
				},
			},
			errFunc: func(t *testing.T, output *AMQP, err error) {
				cfg := output.config
				require.Equal(t, amqp.Table{
					"foo": "bar",
				}, cfg.exchangeArguments)
				require.NoError(t, err)
			},
		},
		{
			name: "username password",
			output: &AMQP{
				Brokers:  []string{"amqp://foo:bar@localhost"},
				Username: config.NewSecret([]byte("telegraf")),
				Password: config.NewSecret([]byte("pa$$word")),
				connect: func(_ *ClientConfig) (Client, error) {
					return NewMockClient(), nil
				},
			},
			errFunc: func(t *testing.T, output *AMQP, err error) {
				cfg := output.config
				require.Equal(t, []amqp.Authentication{
					&amqp.PlainAuth{
						Username: "telegraf",
						Password: "pa$$word",
					},
				}, cfg.auth)

				require.NoError(t, err)
			},
		},
		{
			name: "url support",
			output: &AMQP{
				Brokers: []string{DefaultURL},
				connect: func(_ *ClientConfig) (Client, error) {
					return NewMockClient(), nil
				},
			},
			errFunc: func(t *testing.T, output *AMQP, err error) {
				cfg := output.config
				require.Equal(t, []string{DefaultURL}, cfg.brokers)
				require.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NoError(t, tt.output.Init())
			err := tt.output.Connect()
			tt.errFunc(t, tt.output, err)
		})
	}
}
