package bottom

import (
	"crypto/tls"
	"net/url"
	"strconv"

	"github.com/lrstanley/girc"
)

type Bottom struct {
	Middlewares *Middlewares
	Client      *girc.Client
	ErrorFunc   func(Context, error)
}

func New(user, password, server string, verify bool) (b Bottom, err error) {
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
			InsecureSkipVerify: !verify,
		},
	}

	b.Client = girc.New(config)
	b.Middlewares = NewMiddlewares()
	b.ErrorFunc = b.DefaultErrorFunc

	b.Client.Handlers.Add(girc.PRIVMSG, b.Privmsg)

	return
}

func (b Bottom) Privmsg(_ *girc.Client, e girc.Event) {
	var err error

	ctx := make(Context)
	ctx["sender"] = e.Source.Name
	ctx["message"] = e.Last()

	for _, m := range *b.Middlewares {
		err = m.Do(ctx, e)
		if err != nil {
			b.ErrorFunc(ctx, err)

			return
		}
	}
}

func (b Bottom) DefaultErrorFunc(ctx Context, err error) {
	b.Client.Cmd.Message(ctx["sender"].(string), err.Error())
}
