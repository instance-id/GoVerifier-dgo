package cmdroutes

import (
	. "github.com/instance-id/GoVerifier-dgo/utils"
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

	Dac.SetConfig()

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
