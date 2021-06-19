package bottom

import "github.com/lrstanley/girc"

// Context holds key/value pairs for the current irc event
// like `sender`, or `message`
//
// Because this is a map of strings to interface, you will need
// to typecast whatever comes back
type Context map[string]interface{}

// Middleware provides a way of handling different events in different
// ways, such as gating certain commands
type Middleware interface {
	Do(Context, girc.Event) error
}

// Middlewares holds a set of Middleware implementations, providing
// a set of ways to manipulate insertions
type Middlewares []Middleware

// NewMiddlewares returns an empty set of Middlewares
func NewMiddlewares() *Middlewares {
	m := Middlewares(make([]Middleware, 0))
	return &m
}

// Unshift puts a new Middleware implementation to the front of
// Middlewares
func (m *Middlewares) Unshift(nm Middleware) {
	(*m) = append([]Middleware{nm}, (*m)...)
}

// Push puts a new Middleware implementation to the back of
// Middlewares
func (m *Middlewares) Push(nm Middleware) {
	(*m) = append((*m), nm)
}
