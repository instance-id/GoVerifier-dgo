package cmdroutes

import (
	"fmt"

	. "github.com/instance-id/GoVerifier-dgo/utils"

	"github.com/instance-id/GoVerifier-dgo/models"

	"github.com/sarulabs/di/v2"

	"github.com/Necroforger/dgrouter/exrouter"
)

const DbSetupRoute = "dbsetup"
const DbSetupDescription = "Creates necessary tables in database"

type DbSetup struct {
	di di.Container
}

func (ds *DbSetup) Handle(ctx *exrouter.Context) {
	d := DatabaseAccessContainer(ds.di)

	if !(func() bool { value, _ := d.IsTableExist("verified_users"); return value }() &&
		func() bool { value, _ := d.IsTableExist("user_packages"); return value }() &&
		func() bool { value, _ := d.IsTableExist("asset_packages"); return value }() &&
		func() bool { value, _ := d.IsTableExist("discord_roles"); return value }()) {
		_, err := ctx.Reply("Database schema incomplete. Creating/Updating table schema now...")
		LogFatalf("Verifier had trouble replying: ", err)

		err = d.Sync(new(models.VerifiedUsers), new(models.UserPackages), new(models.AssetPackages), new(models.DiscordRoles))
		if err != nil {
			_, err := ctx.Reply(func() string { result := fmt.Sprintf("Verifier was unable to create tables: %s", err); return result }())
			LogFatalf("Unable to send table creation reply through Discord: ", err)
		}

		resultv, err := d.IsTableExist("verified_users")
		resultp, err := d.IsTableExist("user_packages")
		resulta, err := d.IsTableExist("asset_packages")
		resultr, err := d.IsTableExist("discord_roles")

		_, err = ctx.Reply(func() string {
			return fmt.Sprintf("Schema applied: verified_users: %t - user_packages: %t - asset_packages: %t - discord_roles: %t", resultv, resultp, resulta, resultr)
		}())
		LogFatalf("Unable to send table creation reply through Discord: ", err)

		_, err = ctx.Reply("Database schema creation/update successful")
		LogFatalf("Verifier had trouble replying table creation success: ", err)

	} else {
		_, err := ctx.Reply("Database schema already up to date")
		LogFatalf("Verifier had trouble replying: ", err)
	}
}

func (ds *DbSetup) GetCommand() string {
	return DbSetupRoute
}

func (ds *DbSetup) GetDescription() string {
	return DbSetupDescription
}

func NewDbSetup(di di.Container) *DbSetup {
	return &DbSetup{di: di}
}
