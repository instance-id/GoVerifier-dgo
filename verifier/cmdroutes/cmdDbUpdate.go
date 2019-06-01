package cmdroutes

import (
	"fmt"
	"log"

	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/go-xorm/xorm"
	"github.com/instance-id/GoVerifier-dgo/components"
	"github.com/instance-id/GoVerifier-dgo/models"
	"github.com/sarulabs/di/v2"
	"github.com/sirupsen/logrus"
)

const DbUpdateRoute = "dbupdate"
const DbUpdateDescription = "Creates necessary tables in database"

type DbUpdate struct {
	di di.Container
}

func (e *DbUpdate) Handle(ctx *exrouter.Context) {
	d := e.DataAccessContainer()

	err := d.Sync2(new(models.VerifiedUser), new(models.UserPackages), new(models.AsesetPackages), new(models.DiscordRoles))
	if err != nil {
		_, err := ctx.Reply(func() string { result := fmt.Sprintf("Verifier was unable to create tables: %s", err); return result }())
		if err != nil {
			logrus.Fatalf("Unable to send table creation reply through Discord: %s", err)
		}
	}

	resultv, err := d.IsTableExist("verified_users")
	resultp, err := d.IsTableExist("user_packages")
	resulta, err := d.IsTableExist("asset_packages")
	resultr, err := d.IsTableExist("discord_roles")

	_, err = ctx.Reply(func() string {
		return fmt.Sprintf("Schema applied: verified_users: %t - user_packages: %t - asset_packages: %t - discord_roles: %t", resultv, resultp, resulta, resultr)
	}())
	ErrCheckf("Unable to send table creation reply through Discord:", err)

	_, err = ctx.Reply("Database schema creation/update successful")
	ErrCheckf("Verifier had trouble replying table creation success: ", err)

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

func (e *DbUpdate) DataAccessContainer() *xorm.Engine {
	db, err := e.di.SubContainer()
	if err != nil {
		log.Printf("Error accessing DI container within AddUser module: %s", err)
	}

	dba := db.Get("db").(*components.XormDB).Engine
	return dba
}
