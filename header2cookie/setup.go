package header2cookie

import (
	"strconv"

	"github.com/caddyserver/caddy"
	"github.com/caddyserver/caddy/caddyhttp/httpserver"
)

func init() {
	caddy.RegisterPlugin("header2cookie", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {

	h2c := Header2Cookie{
		CookieName:          "access_token",
		CookieExpireMinutes: 120,
		Expression:          `\s(\w+)$`,
	}

	for c.Next() { // skip the directive name
		if !c.NextArg() { // expect at least one value
			return c.ArgErr() // otherwise it's an error
		}
		value := c.Val() // use the value

		switch value {
		case "CookieName":
			h2c.CookieName = value
		case "CookieExpireMinutes":
			i, err := strconv.Atoi(value)
			if err == nil {
				h2c.CookieExpireMinutes = i
			}
		}
	}

	cfg := httpserver.GetConfig(c)

	mid := func(next httpserver.Handler) httpserver.Handler {
		h2c.Next = next
		return h2c

	}

	cfg.AddMiddleware(mid)

	return nil
}
