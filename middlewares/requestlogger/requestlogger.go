package requestlogger

import (
	"fmt"
	"log"

	"github.com/Tarocch1/kid"
)

var DefaultFormatter = func(c *kid.Ctx) string {
	return fmt.Sprintf("%s %s", c.Method(), c.Url().RequestURI())
}

type Config struct {
	// Skip the middleware when this func return true.
	//
	// Optional. Default: nil
	Skip func(*kid.Ctx) bool

	// Use log.Println rather than fmt.Println.
	//
	// Optional. Default: false
	UseLog bool

	// Formatter formats ctx to log string.
	//
	// Optional. Default: DefaultFormatter
	Formatter func(c *kid.Ctx) string
}

var ConfigDefault = Config{
	Skip:      nil,
	UseLog:    false,
	Formatter: DefaultFormatter,
}

// New creates a new middleware handler
func New(config ...Config) kid.HandlerFunc {
	// Set default config
	cfg := ConfigDefault

	// Override config if provided
	if len(config) > 0 {
		cfg = config[0]

		// Set default values
		if cfg.Formatter == nil {
			cfg.Formatter = ConfigDefault.Formatter
		}
	}

	return func(c *kid.Ctx) error {
		// Don't execute middleware if Skip returns true
		if cfg.Skip != nil && cfg.Skip(c) {
			return c.Next()
		}

		message := cfg.Formatter(c)
		if cfg.UseLog {
			log.Println(message)
		} else {
			fmt.Println(message)
		}

		return c.Next()
	}
}
