package http

import (
	"net/http"
)

var FileServer = http.FileServer
var StripPrefix = http.StripPrefix

type Dir = http.Dir
