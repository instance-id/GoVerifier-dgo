package cmdroutes

import (
	"fmt"
	"log"

	"github.com/sirupsen/logrus"

	"github.com/instance-id/GoVerifier-dgo/models"

	"github.com/instance-id/GoVerifier-dgo/components"

	"github.com/sarulabs/di/v2"

	"github.com/Necroforger/dgrouter/exrouter"
)

const EnsureRoute = "ensure"
const EnsureDescription = "Reloads all actions"

type Ensure struct {
	di di.Container
}

func (e *Ensure) Handle(ctx *exrouter.Context) {

	db, err := e.di.SubContainer()
	if err != nil {
		log.Printf("ERR FROM ENSURE: %s", err)
	}

	d := db.Get("db").(*components.XormDB).Engine
	log.Printf("DATA FROM ENSURE: %s", d.Sync(new(models.ValidatedUsers)))

	results, err := d.Table("validated_users").Exist()
	if err != nil {
		_, err := ctx.Reply("Action!")
		if err != nil {
			log.Printf("Verifier had trouble replying error: %s", err)
		}
	}
	logrus.Infof("Testing MySQL Connection to DB from Verifier: Table 'verified_users' exists? %t", results)

	x, _ := fmt.Printf("Testing MySQL Connection to DB from Verifier: Table 'verified_users' exists? %t", results)
	_, err = ctx.Reply(x)
	if err != nil {
		log.Printf("Verifier had trouble reloading actions: %s", err)
	}
}

func (e *Ensure) GetCommand() string {
	return EnsureRoute
}

func (e *Ensure) GetDescription() string {
	return EnsureDescription
}

func NewEnsure(di di.Container) *Ensure {
	return &Ensure{di: di}
}
