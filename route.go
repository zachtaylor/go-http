package http

type Route interface {
	Match(string) bool
	Responder
}

func Map(s string, f ResponderFunc) {
	router = append(router, &LiteralRoute{s, f})
}
