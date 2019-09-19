package router_test

import (
	"net/http/httptest"
	"testing"

	"ztaylor.me/http/router"
)

func TestPathStarts(t *testing.T) {
	router := router.PathStarts("/hello/")

	r := httptest.NewRequest("", "/hello/", nil)

	if !router.Route(r) {
		t.Log("router path starts /hello/ matches /hello/")
		t.Fail()
	}

	r.URL.Path = "/hello"

	if router.Route(r) {
		t.Log("router path starts /hello/ matches /hello")
		t.Fail()
	}

	r.URL.Path = "/hello/world"

	if !router.Route(r) {
		t.Log("router path starts /hello matches /hello/world")
		t.Fail()
	}
}
