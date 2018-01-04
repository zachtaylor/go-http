package http

type LiteralRoute struct {
	route
	Path string
}

func NewRouteLiteral(path string, r ResponderFunc) Route {
	return &LiteralRoute{
		route: route{r},
		Path:  path,
	}
}

func (route *LiteralRoute) Match(s string) bool {
	return route.Path == s
}
