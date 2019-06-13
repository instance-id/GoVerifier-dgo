package cmdroutes

import (
	. "github.com/instance-id/GoVerifier-dgo/utils"
	"github.com/sarulabs/di/v2"

	"github.com/Necroforger/dgrouter/exrouter"
)

type ExampleSubRoute struct {
	di di.Container
}

func (u *ExampleSubRoute) GetSubRoutes() []Route {
	return []Route{
		NewUser(u.di),
	}
}

func (u *ExampleSubRoute) GetDescription() string {
	return "Testing subroutines"
}

func (u *ExampleSubRoute) Handle(ctx *exrouter.Context) {
	_, err := ctx.Reply("This is a sub route. " + ctx.Msg.Author.Username)
	LogFatalf("Something went wrong: ", err)
}

func (u *ExampleSubRoute) GetCommand() string {
	return "Subroutine"
}

func NewSubRoute(di di.Container) *ExampleSubRoute {
	return &ExampleSubRoute{di}
}
