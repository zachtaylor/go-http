# ztaylor.me/http/...

This package serves as the module root for related golang http libraries, and is not a replacement for `net/http`

Version control is hosted on my private server, and mirrored on Github

## Package `cookies`

Syntactic sugar to read/write cookie values

## Package `handlers`

Additional `http.Handler` types

## Package `json`

Easy encoding/decoding JSON objects

## Package `mux`

Better request routing

## Package `mux/git`

Server routes for git using `github.com/AaronO/go-git-http`

## Package `sessions`

Session creation and management using `ztaylor.me/http/cookies`

## Package `socket`

Websocket server types, connection upgrader, optionally depends on `ztaylor.me/http/sessions`

## Package `track`

Track visitors with `net/http`

## Package `websocket`

Provides websocket connection upgrader, and message handler/mux types
