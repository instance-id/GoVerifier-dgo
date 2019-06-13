package models

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	. "github.com/instance-id/GoVerifier-dgo/data"
	. "github.com/instance-id/GoVerifier-dgo/utils"
)

type PermissionDataAccessObject struct{}

type AvailableRoles struct {
	Roles discordgo.Roles
}

var PermissionDAO *PermissionDataAccessObject

func (p *PermissionDataAccessObject) AddRoles(lc LocalContext, asset string) (bool, string) {
	member, err := lc.Ctx.Ses.GuildMember(lc.Di.Discord.Guild, lc.Ctx.Msg.Author.ID)
	hasError := LogErrorRet(fmt.Sprintf("Could not retrieve member: %s", err), err)
	if hasError {
		return hasError, fmt.Sprintf("Could not retrieve member: %s", err)
	}

	primaryRole := Dac.Discord.Roles["Verified"]
	assetRole := Dac.Discord.Roles[asset]

	hasPrimary := CheckRoleExists(member, primaryRole)

	Log.Infof("PrimaryRole: %s, AssetRole: %s, GuildId: %s, AuthorID: %s ", primaryRole, assetRole, lc.Di.Discord.Guild, lc.Ctx.Msg.Author.ID)

	if !hasPrimary {
		err := lc.Ctx.Ses.GuildMemberRoleAdd(lc.Di.Discord.Guild, lc.Ctx.Msg.Author.ID, primaryRole)
		LogErrorf(fmt.Sprintf("Could not add primary role to user: %s : %s", lc.Ctx.Msg.Author.Username, err), err)
		hasError := LogErrorRet(fmt.Sprintf("Could not add primary role to user: %s : %s", lc.Ctx.Msg.Author.Username, err), err)
		if hasError {
			return hasError, fmt.Sprintf("Could not add primary role to user: %s", lc.Ctx.Msg.Author.Username)
		}
	}

	err = lc.Ctx.Ses.GuildMemberRoleAdd(lc.Di.Discord.Guild, lc.Ctx.Msg.Author.ID, assetRole)
	LogErrorRet(fmt.Sprintf("Could not add asset role to user: %s : %s", lc.Ctx.Msg.Author.ID, err), err)
	hasError = LogErrorRet(fmt.Sprintf("Could not add asset role to user: %s : %s", lc.Ctx.Msg.Author.ID, err), err)
	if hasError {
		return hasError, fmt.Sprintf("Could not add asset role to user: %s", lc.Ctx.Msg.Author.Username)
	}

	return false, fmt.Sprintf("Permissions have been successfully applied to user: %s", lc.Ctx.Msg.Author.Username)
}

func CheckRoleExists(member *discordgo.Member, role string) bool {
	return func() bool {
		for _, v := range member.Roles {
			if v == role {
				return true
			}
		}
		return false
	}()
}
