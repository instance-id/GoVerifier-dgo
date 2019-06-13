package cmdroutes

import (
	"github.com/Necroforger/dgrouter/exrouter"
	. "github.com/instance-id/GoVerifier-dgo/utils"
)

const pingRoute = "ping"
const pingDescription = "responds with pong"

type Ping struct{}

func (p *Ping) Handle(ctx *exrouter.Context) {
	_, err := ctx.Reply("pong")
	LogFatalf("Something went wrong when handling Ping request: ", err)
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
