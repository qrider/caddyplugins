package header2cookie

import (
	"io"

	"net/http"

	"net/http/httptest"

	"testing"

	"github.com/caddyserver/caddy/caddyhttp/httpserver"
)

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {

	// A very simple health check.

	w.WriteHeader(http.StatusOK)

	w.Header().Set("Content-Type", "application/json")

	// In the future we could report back on the status of our DB, or our cache

	// (e.g. Redis) by performing a simple PING, and include them in the response.

	io.WriteString(w, `{"alive": true}`)

}

func TestHandleError(t *testing.T) {

	for _, h2c := range []Header2Cookie{

		Header2Cookie{

			Next: httpserver.HandlerFunc(func(w http.ResponseWriter, r *http.Request) (int, error) {

				r.Header.Add("Authorization", "bearer blah")

				return http.StatusOK, nil

			}),

			CookieName: "access_token",

			CookieExpireMinutes: 120,

			Expression: `\s(\w+)$`,
		},

		Header2Cookie{

			Next: httpserver.HandlerFunc(func(w http.ResponseWriter, r *http.Request) (int, error) {

				r.Header.Add("Authorization", "Bearer blah")

				return http.StatusOK, nil

			}),

			CookieName: "access_token",

			CookieExpireMinutes: 120,

			Expression: `\s(\w+)$`,
		},
	} {

		// setup fake request and response recorder

		rr := httptest.NewRecorder()

		req := httptest.NewRequest(http.MethodGet, "/", nil)

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// A very simple health check.

			w.WriteHeader(http.StatusOK)

			w.Header().Set("Content-Type", "application/json")

			// In the future we could report back on the status of our DB, or our cache

			// (e.g. Redis) by performing a simple PING, and include them in the response.

			io.WriteString(w, `{"`+h2c.CookieName+`": true}`)

		})

		handler.ServeHTTP(rr, req)

		recordedResponse := rr.Result()

		cookies := recordedResponse.Cookies()

		if len(cookies) > 0 && cookies[0].Name != h2c.CookieName {

			t.Log(cookies[0].Name)

			t.Errorf("expected cookie of %s, got %s", string(h2c.CookieName), string(cookies[0].Name))

		}

	}

}
