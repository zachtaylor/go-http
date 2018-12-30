# ztaylor.me/http/...

This serves as the root for related golang http utilities. This package is not intended to be used as a replacement for http.

Version control for this package is hosted on my private source control server and mirrored on Github. Import paths using the Github links will not resolve because I have placed package import comments.

## Package `cookies`

Syntactic sugar to read/write cookie values

## Package `handlers`

Provides additional `http.Handler` types

## Package `mux`

Provides type `Matcher`, used to sort `*http.Request`

Provides additional `Matcher` types

## Package `mux/git`

Provides git routes using `ztaylor.me/http/mux` and `github.com/AaronO/go-git-http`

## Package `sessions`

Provides session creation and management using `ztaylor.me/http/cookies`

## Package `track`

Provides syntactic sugar and interfaces for tracking visitors by IP using `net/http`

## Package `ws`

Provides websocket connection integrations
