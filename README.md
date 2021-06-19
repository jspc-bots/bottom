# bottom

[![codecov](https://codecov.io/gh/jspc/bottom/branch/master/graph/badge.svg)](https://codecov.io/gh/jspc/bottom)
[![GoDoc](https://img.shields.io/badge/pkg.go.dev-doc-blue)](http://pkg.go.dev/github.com/jspc/bottom)
[![Go Report Card](https://goreportcard.com/badge/github.com/jspc/bottom)](https://goreportcard.com/report/github.com/jspc/bottom)

package bottom is an IRC Bot mini-framework (for want of a better word)

It sets up different message handlers for different messages, has the concept of
middlewares, and exposes everything underneath for when you need just a little
more control.

## It's pretty intuitive, but the examples directory will get you where you need to go

It also makes a series of assumptions, such as SASL, and so may not be right for
lots of projects.

## Types

### type [Bottom](/bottom.go#L14)

`type Bottom struct { ... }`

Bottom is the transport logic for an IRC Bot

It handles things like routing messages and error handling

### type [Context](/middleware.go#L10)

`type Context map[string]interface{ ... }`

Context holds key/value pairs for the current irc event
like `sender`, or `message`

Because this is a map of strings to interface, you will need
to typecast whatever comes back

### type [Middleware](/middleware.go#L14)

`type Middleware interface { ... }`

Middleware provides a way of handling different events in different
ways, such as gating certain commands

### type [Middlewares](/middleware.go#L20)

`type Middlewares []Middleware`

Middlewares holds a set of Middleware implementations, providing
a set of ways to manipulate insertions

### type [Router](/router.go#L17)

`type Router struct { ... }`

Router is a Middleware implementation, containing routing logic
for different Events.

These events are pattern matched to a specific RouterFunc

### type [RouterFunc](/router.go#L11)

`type RouterFunc func(sender string, groups []string) error`

RouterFunc is used in routing events

## Sub Packages

* [examples](./examples)

---
Readme created from Go doc with [goreadme](https://github.com/posener/goreadme)
