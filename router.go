package bottom

import (
	"regexp"
	"time"

	"github.com/lrstanley/girc"
)

type RouterFunc func(sender string, groups []string) error

type Router struct {
	routes map[*regexp.Regexp]RouterFunc
}

func NewRouter() *Router {
	return &Router{
		routes: make(map[*regexp.Regexp]RouterFunc),
	}
}

func (r *Router) AddRoute(pattern string, f RouterFunc) (err error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return
	}

	r.routes[re] = f

	return
}

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
