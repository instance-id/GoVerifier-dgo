package cmdroutes

import (
	"bytes"

	. "github.com/instance-id/GoVerifier-dgo/utils"

	"github.com/instance-id/GoVerifier-dgo/appconfig"

	"github.com/sarulabs/di/v2"

	"github.com/bwmarrin/discordgo"

	"github.com/olekukonko/tablewriter"

	"github.com/Necroforger/dgrouter/exrouter"
)

const listrolesRoute = "listroles"
const listrolesDescription = "List all roles on server"

type ListRoles struct {
	di di.Container
}

func (p *ListRoles) Handle(ctx *exrouter.Context) {

	guildRoles, err := ctx.Ses.GuildRoles(p.di.Get("configData").(*appconfig.MainSettings).Discord.Guild)
	LogFatalf("Could not get list of current roles: ", err)

	roleTable := p.renderMarkDownTable(guildRoles)
	_, err = ctx.Reply("```" + roleTable + "```")
	LogFatalf("Something went wrong when handling listroles request: ", err)
}

func (p *ListRoles) GetCommand() string {
	return listrolesRoute
}

func (p *ListRoles) GetDescription() string {
	return listrolesDescription
}

func (p ListRoles) renderMarkDownTable(guildroles discordgo.Roles) string {
	var tableData [][]string

	for _, v := range guildroles {
		row := []string{v.Name, v.ID}
		tableData = append(tableData, row)
	}

	buffer := new(bytes.Buffer)

	table := tablewriter.NewWriter(buffer)
	table.SetHeader([]string{"Role Name:", "Role ID:"})
	table.SetColWidth(40)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.AppendBulk(tableData)
	table.Render()

	return buffer.String()
}

func NewListRoles(di di.Container) *ListRoles {
	return &ListRoles{di: di}
}
