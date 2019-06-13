package data

import (
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bwmarrin/discordgo"
	"github.com/instance-id/GoVerifier-dgo/appconfig"
)

type LocalContext struct {
	Ctx *exrouter.Context
	C   *discordgo.Channel
	Di  *appconfig.MainSettings
}
