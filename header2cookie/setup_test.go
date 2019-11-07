package header2cookie

import (
	"testing"

	"github.com/caddyserver/caddy"

	"github.com/caddyserver/caddy/caddyhttp/httpserver"
)

func TestSetup(t *testing.T) {

	// create test controller and call setup()

	c := caddy.NewTestController("http", "apitoken")

	err := setup(c)

	if err != nil {

		t.Errorf("Expected no errors, but got: %v", err)

	}

	// check that middleware was registered

	m := httpserver.GetConfig(c).Middleware()

	if len(m) == 0 {

		t.Fatal("Expected middleware, but had 0 registered.")

	}

	// check that ApiToken is the registered middleware

	handler := m[0](httpserver.EmptyNext)

	myHandler, ok := handler.(Header2Cookie)

	if !ok {

		t.Fatalf("Expected handler to be type ApiToken, got: %#v", handler)

	}

	// check that the 'next' property is setup

	if !httpserver.SameNext(myHandler.Next, httpserver.EmptyNext) {

		t.Error("'next' field of handler was not set properly")

	}

}
