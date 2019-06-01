package cmdroutes

import (
	"log"

	"github.com/instance-id/GoVerifier-dgo/models"

	"github.com/go-xorm/xorm"
	"github.com/instance-id/GoVerifier-dgo/components"

	"github.com/sarulabs/di/v2"

	"github.com/Necroforger/dgrouter/exrouter"
)

const addUserRoute = "addUser"
const addUserDescription = "Test route to add new user"

type AddUser struct {
	di di.Container
}

func (a *AddUser) Handle(ctx *exrouter.Context) {
	user := models.NewVerifiedUser("MostHated", "M374llic4@gmail.com")
	models.VerifiedUserDAO.AddUser(user, a.di)

	_, err := ctx.Reply("User has been added")
	ErrCheckf("Something went wrong when handling AddUser request: ", err)
}

func (a *AddUser) GetCommand() string {
	return addUserRoute
}

func (a *AddUser) GetDescription() string {
	return addUserDescription
}

func NewAddUser(di di.Container) *AddUser {
	return &AddUser{di: di}
}

func (a *AddUser) DataAccessContainer() *xorm.Engine {
	db, err := a.di.SubContainer()
	if err != nil {
		log.Printf("Error accessing DI container within AddUser module: %s", err)
	}

	dba := db.Get("db").(*components.XormDB).Engine
	return dba
}
