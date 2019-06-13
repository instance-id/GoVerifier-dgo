package cmdroutes

import (
	"bytes"
	"fmt"

	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/instance-id/GoVerifier-dgo/appconfig"
	. "github.com/instance-id/GoVerifier-dgo/data"
	"github.com/instance-id/GoVerifier-dgo/models"
	. "github.com/instance-id/GoVerifier-dgo/utils"
	"github.com/olekukonko/tablewriter"
	"github.com/sarulabs/di/v2"
)

type Verify struct {
	di di.Container
}

const verifyCommand = "verify"
const verifyDescription = "Automatic asset invoice verification"

// --- Setup methods --------------------------------------------------------------------------------------------------------------------------------
func (d *Verify) GetCommand() string {
	return verifyCommand
}

func (d *Verify) GetDescription() string {
	return verifyDescription
}

func NewVerify(di di.Container) *Verify {
	return &Verify{di: di}
}

// --- Module execution : Take input from user (invoice), verify purchase against Unity Asset Store API and apply Discord permission ----------------
func (d *Verify) Handle(ctx *exrouter.Context) {
	var invoiceNum string

	data := DataAccessContainer(d.di)
	c, err := ctx.Ses.UserChannelCreate(ctx.Msg.Author.ID)
	LogFatalf("Could not create direct channel to user: ", err)

	lc := LocalContext{
		Ctx: ctx,
		C:   c,
		Di:  data,
	}

	prompt := "Hello! Which asset would you like to verify? Please reply to be with it's corresponding identifier\n"
	// --- Asset table creation and display -------------------------------------------------------
	table := d.renderMarkDownTable(data)
	assetTable := fmt.Sprintf("\n" + "```" + table + "```")
	prompt += assetTable
	LogFatalf("Something went wrong when handling verify request: ", err)

	// --- Step 1 : Obtain asset to verify  -------------------------------------------------------
	assetChoice := EnterAsset(lc, prompt)

	// --- Step 2 : If asset is found, request invoice for asset. Else: prompt to try again  ------
	if val, ok := data.Assets.Packages[assetChoice]; ok {
		invoiceNum = EnterInvoice(lc, val, assetChoice)
	} else {
		_, err = ctx.Ses.ChannelMessageSend(c.ID, fmt.Sprintf("%s is not a valid choice. Please try again by typing !cmd verify", assetChoice))
		LogFatalf("Could not create direct channel to user: ", err)
		return
	}

	// --- Step 3 : Check if invoice number is already assigned to a user -------------------------
	invoiceResult := models.InvoiceDAO.CheckInvoice(invoiceNum)
	if invoiceResult {
		_, err = ctx.Ses.ChannelMessageSend(c.ID, fmt.Sprintf("Invoice #%s is already in use by another user. Please contact support", assetChoice))
		LogFatalf("Could not create direct channel to user: ", err)
		return
	}

	// --- Step 4 : Obtain email address ----------------------------------------------------------
	emailAddress := EnterEmail(lc)

	// --- Step 5 : Validate against Unity Asset Store API ----------------------------------------
	verified, msg, assetData := models.InvoiceDAO.VerifyInvoice(invoiceNum, assetChoice)
	if !verified {
		_, err := ctx.Ses.ChannelMessageSend(c.ID, msg)
		LogErrorf("Error sending verification results reply to user: ", err)
		return
	}

	// --- Step 6 : Create User Objects -----------------------------------------------------------
	user := models.NewVerifiedUser(ctx.Msg.Author.Username+"#"+ctx.Msg.Author.Discriminator, emailAddress)
	packages := models.NewUserPackages(user, assetData.Invoice, assetChoice, assetData.Date)

	Log.Debugf("USER: %s : %s : %s", user.Id, user.Username, user.Email)
	Log.Debugf("Packages: %s : %s : %s", packages.Id, packages.Username, packages.Purdate)

	// --- Step 7 : Apply permission to channels  -------------------------------------------------
	hadError, msg := models.PermissionDAO.AddRoles(lc, packages.Packages)
	_, err = ctx.Ses.ChannelMessageSend(c.ID, msg)
	LogErrorf("Error adding permission to user: ", err)
	if hadError == true {
		return
	}

	// --- Step 8 : Write database entries --------------------------------------------------------
	hadError, msg = models.InvoiceDAO.AddInvoice(user, packages)
	_, err = ctx.Ses.ChannelMessageSend(c.ID, msg)
	LogErrorf("Error sending verification results reply to user: ", err)
	if hadError == true {
		return
	}

	// --- Step 9 : Completed verification message ------------------------------------------------
	msg = fmt.Sprintf("You have successfully completed verification of %s! Thanks for using Verifier", assetChoice)
	_, err = ctx.Ses.ChannelMessageSend(c.ID, msg)
	LogErrorf("Error adding permission to user: ", err)
	if hadError == true {
		return
	}
}

// --- Asset list table -----------------------------------------------------------------------------------------------------------------------------
func (d *Verify) renderMarkDownTable(data *appconfig.MainSettings) string {
	var tableData [][]string

	// --- Iterates over available asset packages -------------------
	for k, v := range data.Assets.Packages {
		row := []string{k, v}
		tableData = append(tableData, row)
	}

	buffer := new(bytes.Buffer)
	// --- Writes asset data to formatted table to send to user -----
	table := tablewriter.NewWriter(buffer)
	table.SetHeader([]string{"Asset Identifier:", "Asset Name:"})
	table.SetColWidth(60)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.AppendBulk(tableData)
	table.Render()

	return buffer.String()
}
