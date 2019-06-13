package cmdroutes

import (
	. "github.com/instance-id/GoVerifier-dgo/utils"
	"github.com/sarulabs/di/v2"

	"github.com/Necroforger/dgrouter"
	"github.com/Necroforger/dgrouter/exrouter"
)

type User struct {
	di di.Container
}

func (u *User) Register(router *exrouter.Route) *exrouter.Route {
	return router.OnMatch(u.GetCommand(), dgrouter.NewRegexMatcher("user(name)?"), u.Handle)
}

func (u *User) Handle(ctx *exrouter.Context) {
	_, err := ctx.Reply("Your username is " + ctx.Msg.Author.Username)
	LogFatalf("Something went wrong: ", err)
}

func (u *User) GetCommand() string {
	return "username"
}

func (u *User) GetDescription() string {
	return "returns the users username"
}

func NewUser(di di.Container) *User {
	return &User{di: di}
}
