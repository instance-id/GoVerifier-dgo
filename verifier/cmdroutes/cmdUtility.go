package cmdroutes

import (
	"log"

	"github.com/Necroforger/dgrouter/exrouter"
)

const cmdReloadRoute = "reload"
const cmdReloadDescription = "Reloads all actions"

type Reload struct{}

func (p *Reload) Handle(ctx *exrouter.Context) {
	_, err := ctx.Reply("Action!")
	if err != nil {
		log.Print("Verifier had trouble reloading actions: ", err)
	}
}

func (p *Reload) GetCommand() string {
	return cmdReloadRoute
}

func (p *Reload) GetDescription() string {
	return cmdReloadDescription
}

func NewReload() *Reload {
	return &Reload{}
}
