package recovery

import (
	"fmt"

	"github.com/Tarocch1/kid"
)

type Config struct {
	// Optional. Default: nil
	Skip func(*kid.Ctx) bool
}

var ConfigDefault = Config{
	Skip: nil,
}

// New creates a new middleware handler
func New(config ...Config) kid.HandlerFunc {
	// Set default config
	cfg := ConfigDefault

	// Override config if provided
	if len(config) > 0 {
		cfg = config[0]
	}

	// Return new handler
	return func(c *kid.Ctx) (err error) {
		// Don't execute middleware if Skip returns true
		if cfg.Skip != nil && cfg.Skip(c) {
			return c.Next()
		}

		// Catch panics
		defer func() {
			if r := recover(); r != nil {
				var ok bool
				if err, ok = r.(error); !ok {
					// Set error that will call the global error handler
					err = fmt.Errorf("%v", r)
				}
			}
		}()

		// Return err if exist, else move to next handler
		return c.Next()
	}
}
