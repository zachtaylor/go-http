# ztaylor.me/http/...

This package serves as the module root for related golang http libraries, and is not a replacement for `net/http`

Version control is hosted on my private server, and mirrored on Github

## Package `cookies`

Syntactic sugar to read/write cookie values

## Package `handlers`

Additional `http.Handler` types

## Package `mux`

Better request routing

### `mux/acme`

`func Thumbprint` provides server `*Route` for ACME stateless challenge

### `mux/git`

`func Path` provides server `*Route` for git client (git over https using `github.com/AaronO/go-git-http`)

### `mux/goget`

`func Domain` provides server `*Route` for go tool, echo meta data linking go code private domain source

## Package `sessions`

Session creation and management using `ztaylor.me/http/cookies`

## Package `track`

Track visitors with `net/http`

## Package `websocket`

Websocket server types, connection upgrader, uses `ztaylor.me/http/sessions`
