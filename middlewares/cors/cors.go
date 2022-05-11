package cors

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Tarocch1/kid"
)

type Config struct {
	// Skip the middleware when this func return true.
	//
	// Optional. Default: nil
	Skip func(*kid.Ctx) bool

	// Optional. Default value "*"
	AllowOrigins string

	// Optional. Default value "GET,POST,HEAD,PUT,DELETE,PATCH"
	AllowMethods string

	// Optional. Default value "".
	AllowHeaders string

	// Optional. Default value false.
	AllowCredentials bool

	// Optional. Default value "".
	ExposeHeaders string

	// Optional. Default value 0.
	MaxAge int
}

var ConfigDefault = Config{
	Skip:         nil,
	AllowOrigins: "*",
	AllowMethods: strings.Join([]string{
		http.MethodGet,
		http.MethodPost,
		http.MethodHead,
		http.MethodPut,
		http.MethodDelete,
		http.MethodPatch,
	}, ","),
	AllowHeaders:     "",
	AllowCredentials: false,
	ExposeHeaders:    "",
	MaxAge:           0,
}

// New creates a new middleware handler
func New(config ...Config) kid.HandlerFunc {
	// Set default config
	cfg := ConfigDefault

	// Override config if provided
	if len(config) > 0 {
		cfg = config[0]

		// Set default values
		if cfg.AllowMethods == "" {
			cfg.AllowMethods = ConfigDefault.AllowMethods
		}
		if cfg.AllowOrigins == "" {
			cfg.AllowOrigins = ConfigDefault.AllowOrigins
		}
	}

	// Convert string to slice
	allowOrigins := strings.Split(strings.ReplaceAll(cfg.AllowOrigins, " ", ""), ",")

	// Strip white spaces
	allowMethods := strings.ReplaceAll(cfg.AllowMethods, " ", "")
	allowHeaders := strings.ReplaceAll(cfg.AllowHeaders, " ", "")
	exposeHeaders := strings.ReplaceAll(cfg.ExposeHeaders, " ", "")

	// Convert int to string
	maxAge := strconv.Itoa(cfg.MaxAge)

	// Return new handler
	return func(c *kid.Ctx) (err error) {
		// Don't execute middleware if Skip returns true
		if cfg.Skip != nil && cfg.Skip(c) {
			return c.Next()
		}

		// Get origin header
		origin := c.GetHeader("Origin")
		allowOrigin := ""

		// Check allowed origins
		for _, o := range allowOrigins {
			if o == "*" && cfg.AllowCredentials {
				allowOrigin = origin
				break
			}
			if o == "*" || o == origin {
				allowOrigin = o
				break
			}
		}

		// Simple request
		if c.Method() != http.MethodOptions {
			c.AddHeader("Vary", "Origin")
			c.SetHeader("Access-Control-Allow-Origin", allowOrigin)
			if cfg.AllowCredentials {
				c.SetHeader("Access-Control-Allow-Credentials", "true")
			}
			if exposeHeaders != "" {
				c.SetHeader("Access-Control-Expose-Headers", exposeHeaders)
			}
			return c.Next()
		}

		// Preflight request
		c.AddHeader("Vary", "Origin")
		c.AddHeader("Vary", "Access-Control-Request-Method")
		c.AddHeader("Vary", "Access-Control-Request-Headers")
		c.SetHeader("Access-Control-Allow-Origin", allowOrigin)
		c.SetHeader("Access-Control-Allow-Methods", allowMethods)

		// Set Allow-Credentials if set to true
		if cfg.AllowCredentials {
			c.SetHeader("Access-Control-Allow-Credentials", "true")
		}
		if allowHeaders != "" {
			c.SetHeader("Access-Control-Allow-Headers", allowHeaders)
		} else {
			h := c.GetHeader("Access-Control-Request-Headers")
			if h != "" {
				c.SetHeader("Access-Control-Allow-Headers", h)
			}
		}

		// Set MaxAge is set
		if cfg.MaxAge > 0 {
			c.SetHeader("Access-Control-Max-Age", maxAge)
		}

		// Send 204 No Content
		return c.SendStatus(http.StatusNoContent)
	}
}
