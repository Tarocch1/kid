package requestid

import (
	"github.com/Tarocch1/kid"
	"github.com/google/uuid"
)

type Config struct {
	// Skip the middleware when this func return true.
	//
	// Optional. Default: nil
	Skip func(*kid.Ctx) bool

	// Header name that specify the id.
	//
	// Optional. Default: "X-Request-ID"
	Header string

	// Generator generates a new id.
	//
	// Optional. Default: uuid
	Generator func() string
}

var DefaultConfig = Config{
	Skip:      nil,
	Header:    kid.HeaderRequestId,
	Generator: uuid.NewString,
}

// New creates a new middleware handler
func New(config ...Config) kid.HandlerFunc {
	// Set default config
	cfg := DefaultConfig

	// Override config if provided
	if len(config) > 0 {
		cfg = config[0]

		// Set default values
		if cfg.Header == "" {
			cfg.Header = DefaultConfig.Header
		}
		if cfg.Generator == nil {
			cfg.Generator = DefaultConfig.Generator
		}
	}

	return func(c *kid.Ctx) error {
		// Don't execute middleware if Skip returns true
		if cfg.Skip != nil && cfg.Skip(c) {
			return c.Next()
		}

		// Get request id from header if it exits, else generate one.
		rid := c.GetHeader(cfg.Header, cfg.Generator())

		c.Set(kid.CtxRequestId, rid)
		c.SetHeader(cfg.Header, rid)

		return c.Next()
	}
}
