package bottom

import (
	"crypto/tls"
	"net/url"
	"strconv"

	"github.com/lrstanley/girc"
)

// Bottom is the transport logic for an IRC Bot
//
// It handles things like routing messages and error handling
type Bottom struct {
	Middlewares *Middlewares
	Client      *girc.Client
	ErrorFunc   func(Context, error)

	// nick is the name this bot is using
	// it's useful for testing whether messages
	// come from a channel, or directly
	nick string
}

// New accepts some connection parameters, initialises an IRC client,
// and sets up an empty Middlewares list, and a default ErrorFunc
//
// It will error for malformed server addresses. The form it expects is:
//   [irc|ircs]://hostname:port
//
// So:
//   1. `irc://irc.example.com:6667`
//   2. `ircs://irc.example.com:6697`
//
// Are valid, whereas
//
//   1. `irc.example.com`
//
// Is not.
func New(user, password, server string, verifyTLS bool) (b Bottom, err error) {
	u, err := url.Parse(server)
	if err != nil {
		return
	}

	port, err := strconv.Atoi(u.Port())
	if err != nil {
		return
	}

	config := girc.Config{
		Server: u.Hostname(),
		Port:   port,
		Nick:   user,
		User:   user,
		Name:   user,
		SASL: &girc.SASLPlain{
			User: user,
			Pass: password,
		},
		SSL: u.Scheme == "ircs",
		TLSConfig: &tls.Config{
			InsecureSkipVerify: !verifyTLS,
		},
		AllowFlood: true,
	}

	b.Client = girc.New(config)
	b.Middlewares = NewMiddlewares()
	b.ErrorFunc = b.defaultErrorFunc
	b.nick = user

	b.Client.Handlers.Add(girc.PRIVMSG, b.privmsg)

	return
}

func (b Bottom) privmsg(_ *girc.Client, e girc.Event) {
	var err error

	ctx := make(Context)
	ctx["sender"] = e.Source.Name
	ctx["recipient"] = b.nick
	ctx["message"] = e.Last()

	for _, m := range *b.Middlewares {
		err = m.Do(ctx, e)
		if err != nil {
			b.ErrorFunc(ctx, err)

			return
		}
	}
}

func (b Bottom) defaultErrorFunc(ctx Context, err error) {
	b.Client.Cmd.Message(ctx["sender"].(string), err.Error())
}
