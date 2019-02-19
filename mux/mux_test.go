package mux_test

import (
	"testing"

	"ztaylor.me/http/mux"
)

func TestRouterPathStarts(t *testing.T) {
	router := mux.RouterPathStarts("/hello/")

	in := NewInput()

	in.path = "/hello/"

	if !router.Route(in) {
		t.Fail()
	}

	in.path = "/hello/world"

	if !router.Route(in) {
		t.Fail()
	}

	in.path = "/hello"

	if router.Route(in) {
		t.Fail()
	}
}
