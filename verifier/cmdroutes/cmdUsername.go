package cmdroutes

import (
	"log"

	"github.com/sarulabs/di"

	"github.com/Necroforger/dgrouter"
	"github.com/Necroforger/dgrouter/exrouter"
)

type User struct {
	di di.Container
}

func (u *User) Register(router *exrouter.Route) *exrouter.Route {
	return router.OnMatch(u.GetCommand(), dgrouter.NewRegexMatcher("user(name)?"), u.Handle)
}

func (u *User) GetDescription() string {
	return "returns the users username"
}

func (u *User) GetCommand() string {
	return "username"
}

func (u *User) Handle(ctx *exrouter.Context) {
	_, err := ctx.Reply("Your username is " + ctx.Msg.Author.Username)
	if err != nil {
		log.Printf("Something went wrong: %v", err)
	}
}

func NewUser(di di.Container) *User {
	return &User{di: di}
}
