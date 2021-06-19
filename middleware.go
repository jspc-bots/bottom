package bottom

import "github.com/lrstanley/girc"

type Context map[string]interface{}

type Middleware interface {
	Do(Context, girc.Event) error
}

type Middlewares []Middleware

func NewMiddlewares() *Middlewares {
	m := Middlewares(make([]Middleware, 0))
	return &m
}

func (m *Middlewares) Unshift(nm Middleware) {
	(*m) = append([]Middleware{nm}, (*m)...)
}

func (m *Middlewares) Push(nm Middleware) {
	(*m) = append((*m), nm)
}
