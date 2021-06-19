package bottom

import (
	"regexp"
	"time"

	"github.com/lrstanley/girc"
)

// RouterFunc is used in routing events
type RouterFunc func(sender string, groups []string) error

// Router is a Middleware implementation, containing routing logic
// for different Events.
//
// These events are pattern matched to a specific RouterFunc
type Router struct {
	routes map[*regexp.Regexp]RouterFunc
}

// NewRouter returns a new, empty Router, with no routes setup
func NewRouter() *Router {
	return &Router{
		routes: make(map[*regexp.Regexp]RouterFunc),
	}
}

// AddRoute configures a RouterFunc to run when a certain route
// pattern matches.
//
// If many patterns match, the first pattern wins
// If the pattern contains groups, then these are passed as the second arg
// to RouterFunc
func (r *Router) AddRoute(pattern string, f RouterFunc) (err error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return
	}

	r.routes[re] = f

	return
}

// Do implements the Middleware interface.
//
// It matches message contents to route patterns, as passed to AddRoute,
// and calls the associated RouterFunc, passing any regexp groups as it
// goes.
func (r *Router) Do(ctx Context, e girc.Event) (err error) {
	// skip messages older than a minute (assume it's the replayer)
	cutOff := time.Now().Add(0 - time.Minute)
	if e.Timestamp.Before(cutOff) {
		// ignore
		return
	}

	msg := []byte(ctx["message"].(string))
	sender := ctx["sender"].(string)

	for r, f := range r.routes {
		if r.Match(msg) {
			return f(sender, groupsStrings(r.FindAllSubmatch(msg, -1)[0]))
		}
	}

	return
}

func groupsStrings(b [][]byte) (s []string) {
	s = make([]string, len(b))
	for idx, group := range b {
		s[idx] = string(group)
	}

	return
}
