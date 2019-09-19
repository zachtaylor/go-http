package router

import "net/http"

// method satisfies mux.Router by matching the request method
type method string

// Route satisfies Router by matching the given request method
func (router method) Route(r *http.Request) bool {
	return string(router) == r.Method
}

// CONNECT is a Router that returns if Request method is CONNECT
var CONNECT = method("CONNECT")

// DELETE is a Router that returns if Request method is DELETE
var DELETE = method("DELETE")

// GET is a Router that returns if Request method is GET
var GET = method("GET")

// HEAD is a Router that returns if Request method is HEAD
var HEAD = method("HEAD")

// OPTIONS is a Router that returns if Request method is OPTIONS
var OPTIONS = method("OPTIONS")

// POST is a Router that returns if Request method is POST
var POST = method("POST")

// PUT is a Router that returns if Request method is PUT
var PUT = method("PUT")

// TRACE is a Router that returns if Request method is TRACE
var TRACE = method("TRACE")
