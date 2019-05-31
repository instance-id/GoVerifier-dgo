package cmdroutes

import (
	"log"

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
	if err != nil {
		log.Printf("Something went wrong: %v", err)
	}
}

func (u *ExampleSubRoute) GetCommand() string {
	return "Subroutine"
}

func NewSubRoute(di di.Container) *ExampleSubRoute {
	return &ExampleSubRoute{di}
}
