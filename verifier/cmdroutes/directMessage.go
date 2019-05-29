package cmdroutes

import (
	"fmt"
	"log"

	"github.com/Necroforger/dgrouter/exrouter"
)

type DirectMessage struct{}

func (d *DirectMessage) Handle(ctx *exrouter.Context) {
	c, err := ctx.Ses.UserChannelCreate(ctx.Msg.Author.ID)
	if err != nil {
		log.Printf("Could not create direct channel to user: %v", err)
	}

	_, err = ctx.Ses.ChannelMessageSend(c.ID, fmt.Sprintf("This is a direct message to %s", ctx.Msg.Author))
	if err != nil {
		log.Printf("Could not send message: %v", err)
	}
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
