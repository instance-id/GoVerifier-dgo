package cmdroutes

import (
	"github.com/Necroforger/dgrouter/exrouter"
	. "github.com/instance-id/GoVerifier-dgo/utils"
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
	LogFatalf("Something went wrong when handling Avatar request: ", err)
}

func NewAvatar() *Avatar {
	return &Avatar{}
}
