# ztaylor.me/http/...

This package serves as the module root for related golang http libraries, and is not a replacement for `net/http`

Version control is hosted on my private server, and mirrored on Github

## Package `cookies`

Syntactic sugar to read/write cookie values

## Package `handlers`

Additional `http.Handler` types

## Package `mux`

Better request routing

### `mux/acme/Thumbprint`

Server Route for ACME stateless challenge

### `mux/git/Path`

Server Route to provide git over https using `github.com/AaronO/go-git-http`

### `mux/goget/Domain`

Server Route for go tool, provides data to map source onto this domain

## Package `sessions`

Session creation and management using `ztaylor.me/http/cookies`

## Package `track`

Track visitors with `net/http`

## Package `websocket`

Websocket server types, connection upgrader, uses `ztaylor.me/http/sessions`
