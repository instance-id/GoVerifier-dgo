package cmdroutes

import (
	"log"

	"github.com/Necroforger/dgrouter/exrouter"
)

const pingRoute = "ping"
const pingDescription = "responds with pong"

type Ping struct{}

func (p *Ping) Handle(ctx *exrouter.Context) {
	_, err := ctx.Reply("pong")
	if err != nil {
		log.Print("Something went wrong when handling Ping request", err)
	}
}

func (p *Ping) GetCommand() string {
	return pingRoute
}

func (p *Ping) GetDescription() string {
	return pingDescription
}

func NewPing() *Ping {
	return &Ping{}
}
