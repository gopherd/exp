package config

import (
	"context"
	"log/slog"

	"github.com/gopherd/core/typing"
	"github.com/gopherd/exp/spawn"
)

type ClientOptions struct {
	// Source is the configuration source.
	Source string
	// ContentType is the content type of the configuration.
	ContentType ContentType
	// Scopes is the scopes to load.
	Scopes Scopes
	// Name is the namer of the scope: snake_case, camel_case, pascal_case, kebab_case or empty.
	Namer string
	// RefreshInterval is the interval to refresh the configuration.
	RefreshInterval typing.Duration
}

// Client is the configuration client.
//
// Usage:
//
//	type MyConfigHub struct {
//		// ...
//	}
//
//	func NewMyConfigHub() *MyConfigHub {
//		return &MyConfigHub{}
//	}
//
//	func (h *MyConfigHub) Parse(data []byte, decoder encoding.Decoder) error {
//		// ...
//		return nil
//	}
//
//	type MyConfigComponent struct {
//		component.BaseComponent[ClientOptions]
//		*config.Client[*MyConfigHub]
//	}
//
//	func (c *MyConfigComponent) Init(ctx context.Context) error {
//		c.Client = config.NewClient[MyConfigHub](c.Options(), NewMyConfigHub)
//		return c.Client.Init(ctx)
//	}
type Client[H Hub] struct {
	config  *Config[H]
	options ClientOptions
	namer   func(string, string) string
	handle  spawn.Handle
}

// NewClient creates a new configuration client.
func NewClient[H Hub](options ClientOptions, new func() H) *Client[H] {
	return &Client[H]{options: options, config: NewConfig(new)}
}

// Latest returns the latest configuration.
func (c *Client[H]) Latest() H {
	return c.config.Latest()
}

func (c *Client[H]) Init(ctx context.Context) error {
	switch c.options.Namer {
	case "snake_case":
		c.namer = snakeCaseNamer
	case "camel_case":
		c.namer = camelCaseNamer
	case "pascal_case":
		c.namer = pascalCaseNamer
	case "kebab_case":
		c.namer = kebabCaseNamer
	}
	_, err := c.config.Load(ctx, Options{
		Source:      c.options.Source,
		ContentType: c.options.ContentType,
		Scopes:      c.options.Scopes,
		Namer:       c.namer,
	})
	return err
}

func (c *Client[H]) Start(ctx context.Context) error {
	c.handle = spawn.Tick(ctx, c.reload, c.options.RefreshInterval.Value())
	return nil
}

func (c *Client[H]) Shutdown(ctx context.Context) error {
	c.handle.Cancel()
	c.handle.Join(ctx)
	return nil
}

func (c *Client[H]) reload(ctx context.Context) {
	_, err := c.config.Load(ctx, Options{
		Source:      c.options.Source,
		ContentType: c.options.ContentType,
		Scopes:      c.options.Scopes,
		Namer:       c.namer,
	})
	if err != nil {
		slog.Error("failed to load configuration", "error", err)
	}
}
