# bottom

[![codecov](https://codecov.io/gh/jspc/bottom/branch/master/graph/badge.svg)](https://codecov.io/gh/jspc/bottom)
[![GoDoc](https://img.shields.io/badge/pkg.go.dev-doc-blue)](http://pkg.go.dev/github.com/jspc/bottom)
[![Go Report Card](https://goreportcard.com/badge/github.com/jspc/bottom)](https://goreportcard.com/report/github.com/jspc/bottom)

Package bottom is an IRC Bot mini-framework (for want of a better word)

It sets up different message handlers for different messages, has the concept of
middlewares, and exposes everything underneath for when you need just a little
more control.

Everything is pretty intuitive, but the examples directory will get you where you
need to go.

It also makes a series of assumptions, such as SASL, and so may not be right for
lots of projects.

## Types

### type [Bottom](/bottom.go#L14)

`type Bottom struct { ... }`

Bottom is the transport logic for an IRC Bot

It handles things like routing messages and error handling

#### func [New](/bottom.go#L35)

`func New(user, password, server string, verifyTLS bool) (b Bottom, err error)`

New accepts some connection parameters, initialises an IRC client,
and sets up an empty Middlewares list, and a default ErrorFunc

It will error for malformed server addresses. The form it expects is:

```go
[irc|ircs]://hostname:port
```

So:

```go
1. `irc://irc.example.com:6667`
2. `ircs://irc.example.com:6697`
```

Are valid, whereas

```go
1. `irc.example.com`
```

Is not.

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

#### func [NewMiddlewares](/middleware.go#L23)

`func NewMiddlewares() *Middlewares`

NewMiddlewares returns an empty set of Middlewares

#### func (*Middlewares) [Push](/middleware.go#L36)

`func (m *Middlewares) Push(nm Middleware)`

Push puts a new Middleware implementation to the back of
Middlewares

#### func (*Middlewares) [Unshift](/middleware.go#L30)

`func (m *Middlewares) Unshift(nm Middleware)`

Unshift puts a new Middleware implementation to the front of
Middlewares

### type [Router](/router.go#L22)

`type Router struct { ... }`

Router is a Middleware implementation, containing routing logic
for different Events.

These events are pattern matched to a specific RouterFunc

#### func [NewRouter](/router.go#L27)

`func NewRouter() *Router`

NewRouter returns a new, empty Router, with no routes setup

#### func (*Router) [AddRoute](/router.go#L39)

`func (r *Router) AddRoute(pattern string, f RouterFunc) (err error)`

AddRoute configures a RouterFunc to run when a certain route
pattern matches.

If many patterns match, the first pattern wins
If the pattern contains groups, then these are passed as the second arg
to RouterFunc

#### func (*Router) [Do](/router.go#L55)

`func (r *Router) Do(ctx Context, e girc.Event) (err error)`

Do implements the Middleware interface.

It matches message contents to route patterns, as passed to AddRoute,
and calls the associated RouterFunc, passing any regexp groups as it
goes.

### type [RouterFunc](/router.go#L16)

`type RouterFunc func(sender, channel string, groups []string) error`

RouterFunc is used in routing events

Functions should expect the following information:

```go
1. sender - the nick of the author of the message
2. channel - the channel the message was sent in (if sender == channel, then assume a private message)
3. groups - any regexp groups extracted from the message. groups[0] is *always* the full message
```

## Sub Packages

* [examples](./examples)

---
Readme created from Go doc with [goreadme](https://github.com/posener/goreadme)
