package ws

// Matcher tests a Message
type Matcher interface {
	Match(*Message) bool
}

// Handler handles a Message
type Handler interface {
	ServeWS(*Message)
}

// Route holds pointers to Matcher and Handler
//
// Route provides Router
type Route struct {
	Matcher
	Handler
}

// Router provides Matcher and Handler
type Router interface {
	Matcher
	Handler
}
