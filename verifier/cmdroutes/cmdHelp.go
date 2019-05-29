package cmdroutes

import (
	"bytes"
	"fmt"
	"log"

	"github.com/olekukonko/tablewriter"

	"github.com/Necroforger/dgrouter/exrouter"
)

type cmdHelp struct {
	router *exrouter.Route
}

func (h *cmdHelp) Handle(ctx *exrouter.Context) {
	var helpMsg string
	if h.router.Name != "" {
		helpMsg = fmt.Sprintf("## Help: \n\n Below are the usable commands for: %s\n\n ", h.router.Name)
	} else {
		helpMsg = fmt.Sprintln("## Help:\n\n Below are the root commands: ")
	}

	table := h.renderMarkDownTable()

	_, err := ctx.Reply("```" + helpMsg + table + "```")
	if err != nil {
		log.Print("Help info request did not complete properly.")
	}
}

func (h cmdHelp) GetCommand() string {
	return "help"
}

func (h cmdHelp) GetDescription() string {
	return "Prints this help menu"
}

func (h cmdHelp) renderMarkDownTable() string {
	var tableData [][]string

	for _, v := range h.router.Routes {
		row := []string{v.Name, v.Description}
		tableData = append(tableData, row)
	}

	buffer := new(bytes.Buffer)

	table := tablewriter.NewWriter(buffer)
	table.SetHeader([]string{"Command Name:", "Command Description:"})
	table.SetColWidth(40)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.AppendBulk(tableData)
	table.Render()

	return buffer.String()
}

func NewHelpRoute(router *exrouter.Route) *cmdHelp {
	return &cmdHelp{router: router}
}
