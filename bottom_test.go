package bottom

import (
	"fmt"
	"testing"
	"time"

	"github.com/lrstanley/girc"
)

func TestNew(t *testing.T) {
	for _, test := range []struct {
		name        string
		server      string
		expectError bool
	}{
		{"happy path", "ircs://irc.example.com:6697", false},
		{"empty server", "", true},
		{"garbage server", "\b\b////b\b\b\b///:::/", true},
		{"missing port", "irc://irc.example.com", true},
	} {
		t.Run(test.name, func(t *testing.T) {
			b, err := New("", "", test.server, true)
			if err == nil && test.expectError {
				t.Fatalf("expected error")
			}

			if err != nil && !test.expectError {
				t.Fatalf("unexpected error: %+v", err)
			}

			// There's no point continuing- it's all going to be empty
			if test.expectError {
				return
			}

			t.Run("Middlewares", func(t *testing.T) {
				if b.Middlewares == nil {
					t.Errorf("No middlewares set")
				}
			})

			t.Run("Client", func(t *testing.T) {
				if b.Client == nil {
					t.Errorf("No middlewares set")
				}
			})

			t.Run("ErrorFunc", func(t *testing.T) {
				if b.ErrorFunc == nil {
					t.Errorf("No middlewares set")
				}
			})
		})
	}
}

func TestBottom_Privmsg(t *testing.T) {
	var count int

	b, _ := New("", "", "ircs://irc.example.com:6697", true)

	r := NewRouter()
	r.AddRoute("(i?)PATTERN", func(_ string, _ []string) error {
		count++

		return fmt.Errorf("an error")
	})

	b.Middlewares.Push(r)

	if b.Middlewares == nil || len(*b.Middlewares) != 1 {
		t.Errorf("middlewares should exist, and contain one thing: %v", b.Middlewares)
	}

	t.Run("invocation count", func(t *testing.T) {
		b.Privmsg(nil, girc.Event{
			Source:    &girc.Source{Name: "#testing"},
			Command:   "PRIVMSG",
			Params:    []string{"#testing", "PATTERN"},
			Timestamp: time.Now(),
		})

		if count != 1 {
			t.Errorf("expected 1, received %d", count)
		}
	})

	t.Run("too old", func(t *testing.T) {
		count = 0
		b.Privmsg(nil, girc.Event{
			Source:    &girc.Source{Name: "#testing"},
			Command:   "PRIVMSG",
			Params:    []string{"#testing", "PATTERN"},
			Timestamp: time.Now().Add(0 - time.Hour),
		})

		if count > 0 {
			t.Errorf("expected 0, received %d", count)
		}
	})

	t.Run("doesn't match pattern", func(t *testing.T) {
		count = 0
		b.Privmsg(nil, girc.Event{
			Source:    &girc.Source{Name: "#testing"},
			Command:   "PRIVMSG",
			Params:    []string{"#testing", "some other message"},
			Timestamp: time.Now(),
		})

		if count > 0 {
			t.Errorf("expected 0, received %d", count)
		}
	})
}
