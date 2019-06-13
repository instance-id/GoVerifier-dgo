package cmdroutes

import (
	"fmt"

	. "github.com/instance-id/GoVerifier-dgo/utils"

	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/go-xorm/xorm"
	"github.com/instance-id/GoVerifier-dgo/models"
	"github.com/sarulabs/di/v2"
)

const DbUpdateRoute = "dbupdate"
const DbUpdateDescription = "Creates necessary tables in database"

type DbUpdate struct {
	di di.Container
}

func (e *DbUpdate) Handle(ctx *exrouter.Context) {
	err := Dba.Sync2(new(models.VerifiedUser), new(models.UserPackages), new(models.AssetPackages), new(models.DiscordRoles))
	if err != nil {
		_, err := ctx.Reply(func() string { result := fmt.Sprintf("Verifier was unable to create tables: %s", err); return result }())
		LogFatalf("Unable to send table creation reply through Discord: ", err)
	}

	resultv, err := Dba.IsTableExist("verified_users")
	resultp, err := Dba.IsTableExist("user_packages")
	resulta, err := Dba.IsTableExist("asset_packages")
	resultr, err := Dba.IsTableExist("discord_roles")

	_, err = ctx.Reply(func() string {
		return fmt.Sprintf("Schema applied: verified_users: %t - user_packages: %t - asset_packages: %t - discord_roles: %t", resultv, resultp, resulta, resultr)
	}())
	LogFatalf("Unable to send table creation reply through Discord:", err)

	_, err = ctx.Reply("Database schema creation/update successful")
	LogFatalf("Verifier had trouble replying table creation success: ", err)

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

// --- Database Functions ------------------------------------------------
func CheckExists(d *xorm.Engine) bool {

	if func() bool { value, _ := d.IsTableExist("verified_users"); return value }() &&
		func() bool { value, _ := d.IsTableExist("user_packages"); return value }() &&
		func() bool { value, _ := d.IsTableExist("asset_packages"); return value }() &&
		func() bool { value, _ := d.IsTableExist("discord_roles"); return value }() {
		return true
	} else {
		return false
	}
}
