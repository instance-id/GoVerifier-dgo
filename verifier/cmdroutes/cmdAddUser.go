package cmdroutes

import (
	"time"

	"github.com/instance-id/GoVerifier-dgo/models"
	. "github.com/instance-id/GoVerifier-dgo/utils"

	"github.com/sarulabs/di/v2"

	"github.com/Necroforger/dgrouter/exrouter"
)

const addUserRoute = "adduser"
const addUserDescription = "Test route to add new user"

type AddUser struct {
	di di.Container
}

func (a *AddUser) Handle(ctx *exrouter.Context) {
	user := models.NewVerifiedUser(ctx.Msg.Author.Username+"#"+ctx.Msg.Author.Discriminator, "M374llic4@gmail.com")
	packages := models.NewUserPackages(user, "123123123", "SCT", time.Date(2018, 12, 25, 0, 0, 0, 0, time.Local))
	models.InvoiceDAO.AddInvoice(user, packages)

	_, err := ctx.Reply("User and packages have been added")
	LogFatalf("Something went wrong when handling AddUser request: ", err)

	//models.VerifiedUserDAO.AddUser(user, a.di)

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
