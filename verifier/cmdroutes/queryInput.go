package cmdroutes

import (
	"errors"
	"time"

	. "github.com/instance-id/GoVerifier-dgo/utils"

	"github.com/Necroforger/dgrouter/exrouter"

	"github.com/bwmarrin/discordgo"
)

func QueryInput(channel *discordgo.Channel, ctx *exrouter.Context, prompt string, timeout time.Duration) (*discordgo.Message, error) {
	msg, err := ctx.Ses.ChannelMessageSend(channel.ID, prompt)
	if err != nil {
		return nil, err
	}

	if channel.Type != discordgo.ChannelTypeDM {
		defer func() {
			err := ctx.Ses.ChannelMessageDelete(msg.ChannelID, msg.ID)
			LogFatalf("Unable to delete message: ", err)
		}()
	}

	timeoutChan := make(chan int)
	go func() {
		time.Sleep(timeout)
		timeoutChan <- 0
	}()

	for {
		select {
		case userMsg := <-NextMessageCreateC(ctx.Ses):
			if userMsg.Author.ID != ctx.Msg.Author.ID {
				continue
			}

			if channel.Type != discordgo.ChannelTypeDM {
				err := ctx.Ses.ChannelMessageDelete(userMsg.ChannelID, userMsg.ID)
				LogFatalf("Unable to delete message: ", err)
			}
			return userMsg.Message, nil
		case <-timeoutChan:
			return nil, errors.New("Timed out")
		}
	}
}

// NextMessageCreateC returns a channel for the next MessageCreate event
func NextMessageCreateC(s *discordgo.Session) chan *discordgo.MessageCreate {
	out := make(chan *discordgo.MessageCreate)
	s.AddHandlerOnce(func(_ *discordgo.Session, e *discordgo.MessageCreate) {
		out <- e
	})
	return out
}

// NextMessageReactionAddC returns a channel for the next MessageReactionAdd event
func NextMessageReactionAddC(s *discordgo.Session) chan *discordgo.MessageReactionAdd {
	out := make(chan *discordgo.MessageReactionAdd)
	s.AddHandlerOnce(func(_ *discordgo.Session, e *discordgo.MessageReactionAdd) {
		out <- e
	})
	return out
}
