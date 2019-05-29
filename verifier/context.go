package verifier

import (
	"github.com/bwmarrin/discordgo"
)

type Context struct {
	Discord      *discordgo.Session
	Guild        *discordgo.Guild
	VoiceChannel *discordgo.Channel
	TextChannel  *discordgo.Channel
	User         *discordgo.User
	Message      *discordgo.MessageCreate
	Args         []string
	Config       *Config
}

func CreateContext(discord *discordgo.Session,
	guild *discordgo.Guild,
	textChannel *discordgo.Channel,
	user *discordgo.User,
	message *discordgo.MessageCreate,
	config *Config) *Context {

	ctx := new(Context)
	ctx.Discord = discord
	ctx.Guild = guild
	ctx.TextChannel = textChannel
	ctx.User = user
	ctx.Message = message
	ctx.Config = config
	return ctx
}
