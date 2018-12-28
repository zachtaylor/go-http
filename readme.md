# import "ztaylor.me/http/..."

This serves as the root for related golang http utilities. This package is not intended to be used as a replacement for http. Version control for this package is hosted on my private source control server and mirrored on Github. Import paths using the Github links will not resolve because I have placed package import comments.

## Package `cookies`

Syntactic sugar for `net/http`

## Package `handlers`

Provides helpful `net/http.Handler` types

## Package `mux`

Provides a router using `net/http.Handler` and `Matcher`

Routers can be combined to create multi-host environments

## Package `mux/git`

Provides git routes using `ztaylor.me/mux` and `github.com/AaronO/go-git-http`

## Package `sessions`

Provides session creation and management using `ztaylor.me/cookies`

## Package `track`

Provides syntactic sugar and interfaces for tracking visitors by IP using `net/http`

## Package `ws`

Provides websocket connection integrations
