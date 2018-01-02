package http

type Responder interface {
	Respond(Agent, *Request) error
}

type ResponderFunc func(Agent, *Request) error

func (f ResponderFunc) Respond(a Agent, r *Request) error {
	return f(a, r)
}

func NewResponder(f ResponderFunc) Responder {
	return f
}
