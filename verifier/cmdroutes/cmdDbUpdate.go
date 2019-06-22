package cmdroutes

import (
	"fmt"

	. "github.com/instance-id/GoVerifier-dgo/utils"

	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/instance-id/GoVerifier-dgo/models"
	"github.com/sarulabs/di/v2"
)

const DbUpdateRoute = "dbupdate"
const DbUpdateDescription = "Creates necessary tables in database"

type DbUpdate struct {
	di di.Container
}

func (e *DbUpdate) Handle(ctx *exrouter.Context) {
	err := Dba.Sync2(new(models.VerifiedUsers), new(models.UserPackages), new(models.AssetPackages), new(models.DiscordRoles))
	if err != nil {
		_, err := ctx.Reply(func() string { result := fmt.Sprintf("Verifier was unable to create tables: %s", err); return result }())
		LogFatalf("Unable to send table creation reply through Discord: ", err)
	}

	resultv, err := Dba.IsTableExist("verified_users")
	LogFatalf(fmt.Sprintf("Verifier could not update table: verified_users : "), err)
	resultp, err := Dba.IsTableExist("user_packages")
	LogFatalf(fmt.Sprintf("Verifier could not update table: user_packages : "), err)
	resulta, err := Dba.IsTableExist("asset_packages")
	LogFatalf(fmt.Sprintf("Verifier could not update table: asset_packages : "), err)
	resultr, err := Dba.IsTableExist("discord_roles")
	LogFatalf(fmt.Sprintf("Verifier could not update table: discord_roles : "), err)

	if resultv && resultp && resulta && resultr {
		_, err = ctx.Reply(func() string {
			return fmt.Sprintf("Schema applied: verified_users: %t - user_packages: %t - asset_packages: %t - discord_roles: %t", resultv, resultp, resulta, resultr)
		}())
		LogFatalf("Unable to send table creation reply through Discord:", err)
		_, err = ctx.Reply("Database schema creation/update successful")
		LogFatalf("Verifier had trouble replying table creation success: ", err)

	} else {

	}

}

func (e *DbUpdate) GetCommand() string {
	return DbUpdateRoute
}

func (e *DbUpdate) GetDescription() string {
	return DbUpdateDescription
}

func NewDbUpdate(di di.Container) *DbUpdate {
	return &DbUpdate{di: di}
}
