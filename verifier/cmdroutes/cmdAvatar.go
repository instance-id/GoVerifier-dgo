package cmdroutes

import (
	"log"

	"github.com/Necroforger/dgrouter/exrouter"
)

const avatarRoute = "avatar"
const avatarDescription = "returns the user's avatar"

type Avatar struct{}

func (a *Avatar) GetDescription() string {
	return avatarDescription
}

func (a *Avatar) GetCommand() string {
	return avatarRoute
}

func (a *Avatar) Handle(ctx *exrouter.Context) {
	_, err := ctx.Reply(ctx.Msg.Author.AvatarURL("2048"))
	if err != nil {
		log.Print("Something went wrong when handling Avatar request", err)
	}
}

func NewAvatar() *Avatar {
	return &Avatar{}
}
