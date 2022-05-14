package requestlogger

import (
	"fmt"

	"github.com/Tarocch1/kid"
)

var DefaultFormatter = func(c *kid.Ctx) (string, map[string]interface{}) {
	message := fmt.Sprintf("%s %s", c.Method(), c.Url().RequestURI())
	extra := map[string]interface{}{
		"body":   string(c.Body()),
		"header": c.Header(),
	}
	return message, extra
}

type Config struct {
	// Skip the middleware when this func return true.
	//
	// Optional. Default: nil
	Skip func(*kid.Ctx) bool

	// The module for logger.
	//
	// Optional. Default: "HTTP Request"
	Module string

	// Formatter formats ctx to log string.
	//
	// Optional. Default: DefaultFormatter
	Formatter func(c *kid.Ctx) (string, map[string]interface{})
}

var DefaultConfig = Config{
	Skip:      nil,
	Module:    "HTTP Request",
	Formatter: DefaultFormatter,
}

// New creates a new middleware handler
func New(config ...Config) kid.HandlerFunc {
	// Set default config
	cfg := DefaultConfig

	// Override config if provided
	if len(config) > 0 {
		cfg = config[0]

		// Set default values
		if cfg.Module == "" {
			cfg.Module = DefaultConfig.Module
		}
		if cfg.Formatter == nil {
			cfg.Formatter = DefaultConfig.Formatter
		}
	}

	logger := kid.NewLogger(cfg.Module)

	return func(c *kid.Ctx) error {
		// Don't execute middleware if Skip returns true
		if cfg.Skip != nil && cfg.Skip(c) {
			return c.Next()
		}

		message, extra := cfg.Formatter(c)
		logger.Info(c, message, extra)

		return c.Next()
	}
}
