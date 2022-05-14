package basicauth

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/Tarocch1/kid"
)

type Config struct {
	// Skip the middleware when this func return true.
	//
	// Optional. Default: nil
	Skip func(*kid.Ctx) bool

	// Users defines the allowed credentials
	//
	// Required. Default: map[string]string{}
	Users map[string]string

	// Realm is a string to define realm attribute of BasicAuth.
	// the realm identifies the system to authenticate against
	// and can be used by clients to save credentials
	//
	// Optional. Default: "Restricted".
	Realm string

	// Authorizer defines a function you can pass
	// to check the credentials however you want.
	// It will be called with a username and password
	// and is expected to return true or false to indicate
	// that the credentials were approved or not.
	//
	// Optional. Default: nil
	Authorizer func(string, string) bool

	// Unauthorized defines the response body for unauthorized responses.
	// By default it will return with a 401 Unauthorized and the correct WWW-Auth header
	//
	// Optional. Default: nil
	Unauthorized kid.HandlerFunc
}

var DefaultConfig = Config{
	Skip:         nil,
	Users:        map[string]string{},
	Realm:        "Restricted",
	Authorizer:   nil,
	Unauthorized: nil,
}

// New creates a new middleware handler
func New(config ...Config) kid.HandlerFunc {
	// Set default config
	cfg := DefaultConfig

	// Override config if provided
	if len(config) > 0 {
		cfg = config[0]

		// Set default values
		if cfg.Users == nil {
			cfg.Users = DefaultConfig.Users
		}
		if cfg.Realm == "" {
			cfg.Realm = DefaultConfig.Realm
		}
		if cfg.Authorizer == nil {
			cfg.Authorizer = func(username, password string) bool {
				userPwd, exist := cfg.Users[username]
				return exist && password == userPwd
			}
		}
		if cfg.Unauthorized == nil {
			cfg.Unauthorized = func(c *kid.Ctx) error {
				c.SetHeader(kid.HeaderWWWAuthenticate, "basic")
				return c.SendStatus(http.StatusUnauthorized)
			}
		}
	}

	return func(c *kid.Ctx) error {
		// Don't execute middleware if Skip returns true
		if cfg.Skip != nil && cfg.Skip(c) {
			return c.Next()
		}

		// Get authorization header
		auth := c.GetHeader(kid.HeaderAuthorization)

		// Check if the header contains content besides "basic"
		if len(auth) <= 6 || strings.ToLower(auth[:5]) != "basic" {
			return cfg.Unauthorized(c)
		}

		// Decode the header contents
		raw, err := base64.StdEncoding.DecodeString(auth[6:])
		if err != nil {
			return cfg.Unauthorized(c)
		}

		// Get the credentials
		creds := string(raw)

		// Check if the credentials are in the correct form
		// which is "username:password".
		index := strings.Index(creds, ":")
		if index == -1 {
			return cfg.Unauthorized(c)
		}

		// Get the username and password
		username := creds[:index]
		password := creds[index+1:]

		if cfg.Authorizer(username, password) {
			c.Set(kid.CtxBasicAuthUsername, username)
			c.Set(kid.CtxBasicAuthPassword, password)
			return c.Next()
		}

		// Authentication failed
		return cfg.Unauthorized(c)
	}
}
