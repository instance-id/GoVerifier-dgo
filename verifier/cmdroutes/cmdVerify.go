package cmdroutes

import (
	"fmt"
	"log"

	"github.com/sarulabs/di/v2"

	"github.com/Necroforger/dgrouter/exrouter"
)

type Verify struct {
	di di.Container
}

const verifyCommand = "verify"
const verifyDescription = "Automatic asset invoice verification"

func (d *Verify) Handle(ctx *exrouter.Context) {
	c, err := ctx.Ses.UserChannelCreate(ctx.Msg.Author.ID)
	if err != nil {
		log.Printf("Could not create direct channel to user: %v", err)
	}

	_, err = ctx.Ses.ChannelMessageSend(c.ID, fmt.Sprintf("Hello! Which asset would you like to verify? Please reply to be with it's cooresponding number."))
	if err != nil {
		log.Printf("Could not send message: %v", err)
	}

	//assetTable := p.renderMarkDownTable(guildRoles)
	//_, err = ctx.Reply("```" + assetTable + "```")
	//if err != nil {
	//	log.Print("Something went wrong when handling listroles request", err)
	//}
}

func (d *Verify) GetCommand() string {
	return verifyCommand
}

func (d *Verify) GetDescription() string {
	return verifyDescription
}

func NewVerify(di di.Container) *Verify {
	return &Verify{di: di}
}
