package header2cookie

import (
	"net/http"

	"regexp"

	"time"

	"github.com/caddyserver/caddy/caddyhttp/httpserver"
)

//Header2Cookie translates authorization header value to cookie value
type Header2Cookie struct {
	Expression          string
	CookieName          string
	CookieExpireMinutes int
	Next                httpserver.Handler
}

func (h Header2Cookie) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {

	auth := r.Header.Get("Authorization")

	//rg, _ := regexp.Compile(`\s(\w+)$`)
	rg, _ := regexp.Compile(h.Expression)

	token := rg.FindString(auth)

	if token != "" {

		cur := time.Now()

		expire := cur.Add(time.Minute * time.Duration(h.CookieExpireMinutes))

		addCookie(w, h.CookieName, token, expire)

	}

	return h.Next.ServeHTTP(w, r)

}

func addCookie(w http.ResponseWriter, name string, value string, expire time.Time) {

	cookie := http.Cookie{

		Name: name,

		Value: value,

		Expires: expire,
	}

	http.SetCookie(w, &cookie)

}
