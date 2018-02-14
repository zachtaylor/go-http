package http

import (
	"net/http"
)

var FileServer = http.FileServer

type Dir http.Dir
