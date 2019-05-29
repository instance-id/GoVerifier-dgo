package verifier

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type UserRoles struct {
	Roles []*discordgo.Role
}

// GetRole returns UserRoles struct pointer
func (ctx *Context) GetRoles() *UserRoles {
	var userRoles = new(UserRoles)
	member, err := ctx.Discord.GuildMember(Guild.ID, User.ID)
	if err != nil {
		fmt.Println("Getting member error: " + err.Error())
	}
	for _, grole := range Guild.Roles {
		for _, urole := range member.Roles {
			if grole.ID == urole {
				userRoles.Roles = append(userRoles.Roles, grole)
			}
		}
	}
	return userRoles
}

// ExistsName checks if user role nema exists on user
func (r *UserRoles) ExistsName(name string) bool {
	for _, val := range r.Roles {
		if val.Name == name {
			return true
		}
	}
	return false
}
