package cmdroutes

import (
	"github.com/sarulabs/di/v2"

	"github.com/Necroforger/dgrouter/exrouter"
)

const TestRoute = "test"
const TestDescription = "Test Route"

type Test struct {
	di di.Container
}

func (t *Test) Handle(ctx *exrouter.Context) {
	//c, err := ctx.Ses.UserChannelCreate(ctx.Msg.Author.ID)
	//LogFatalf("Could not create direct channel to user: ", err)

	//lc := data.LocalContext{
	//	Ctx: ctx,
	//	C:   c,
	//	Di:  Dac,
	//}

	////_, err = Dac.SetConfig()
	//LogFatalf("Could not write to config file: ", err)
	//if err != nil {
	//	msg := fmt.Sprintf("Could not write to config file: %s", err)
	//	_, err = ctx.Ses.ChannelMessageSend(c.ID, msg)
	//	LogFatalf("Could not send reply: ", err)
	//}
	//
	//db, err := t.di.SubContainer()
	//LogFatalf("Error accessing DI container within AddUser module: ", err)
	//
	////database := db.Get("dbConn").(*components.DbConfig).Db
	//
	////_, err = database.SetDbConfig()
	//LogFatalf("Could not write to config file: ", err)
	//if err != nil {
	//	msg := fmt.Sprintf("Could not write to dbconfig file: %s", err)
	//	_, err = ctx.Ses.ChannelMessageSend(c.ID, msg)
	//	LogFatalf("Could not send reply: ", err)
	//}

}

func (t *Test) GetCommand() string {
	return TestRoute
}

func (t *Test) GetDescription() string {
	return TestDescription
}

func NewTest(di di.Container) *Test {
	return &Test{di: di}
}
