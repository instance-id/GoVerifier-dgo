package cmdroutes

import (
	"fmt"

	. "github.com/instance-id/GoVerifier-dgo/utils"

	"github.com/Necroforger/dgrouter/exrouter"
)

type DirectMessage struct{}

func (d *DirectMessage) Handle(ctx *exrouter.Context) {
	c, err := ctx.Ses.UserChannelCreate(ctx.Msg.Author.ID)
	LogFatalf("Could not create direct channel to user: ", err)

	_, err = ctx.Ses.ChannelMessageSend(c.ID, fmt.Sprintf("This is a direct message to %s", ctx.Msg.Author))
	LogFatalf("Could not send message: ", err)
}

func (d *DirectMessage) GetCommand() string {
	return "dm"
}

func (d *DirectMessage) GetDescription() string {
	return "Receive a direct message from the bot"
}

func NewDirectMessage() *DirectMessage {
	return &DirectMessage{}
}
