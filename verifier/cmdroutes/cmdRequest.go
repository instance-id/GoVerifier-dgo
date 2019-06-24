// --- Test route for verification process --------------------------------------------------------------------------------- TODO Remove test route
package cmdroutes

import (
	"time"

	"github.com/Necroforger/dgrouter/exrouter"
	. "github.com/instance-id/GoVerifier-dgo/utils"
	"github.com/parnurzeal/gorequest"
	"github.com/sarulabs/di/v2"
	"github.com/tidwall/gjson"
)

const RequestRoute = "request"
const RequestDescription = "Creates web request"

type Invoices struct {
	PriceExvat   string    `json:"price_exvat"`
	Downloaded   string    `json:"downloaded"`
	Date         time.Time `json:"date"`
	Quantity     string    `json:"quantity"`
	Reason       string    `json:"reason"`
	OtherLicense string    `json:"other_license"`
	Package      string    `json:"package"`
	Currency     string    `json:"currency"`
	Refunded     string    `json:"refunded"`
	Invoice      string    `json:"invoice"`
}

type Request struct {
	di di.Container
}

func (r *Request) GetCommand() string {
	return RequestRoute
}

func (r *Request) GetDescription() string {
	return RequestDescription
}

func NewRequest(di di.Container) *Request {
	return &Request{di: di}
}

func (r *Request) Handle(ctx *exrouter.Context) {
	assetChoice := "UFPS1"
	assetName := Dac.Assets.Packages[assetChoice]
	LogInfof("Asset Name: ", assetName)

	apiKey := Dac.Assets.ApiKeys[assetChoice]
	invoiceNum := "123123123"

	url := "https://api.assetstore.unity3d.com:443/publisher/v1/invoice/verify.json"

	request := gorequest.New()
	resp, body, err := request.Set("Accept", "application/json").Get(url).Query("key=" + apiKey + "&" + "invoice=" + invoiceNum).End()
	if err != nil {
		Log.Warnf("Received warning from Server: %s", err)
	}
	if len(body) <= 0 {
		Log.Warnf("Received empty reply from server. Status code received: %s", resp.StatusCode)
		return
	}

	assetData := gjson.Get(body, "invoices.#[package=="+assetName+"]").Map()
	Log.Infof("Data %v", assetData)

	if len(assetData) <= 1 {
		Log.Warnf("Could not verify: %s for invoice %s. Please check the invoice number and try again.", assetName, invoiceNum)
		return
	}

	convertedDate := ConvertTime(assetData["date"].String())
	Log.Infof("Converted Time: %v", convertedDate)

	asset := Invoices{
		PriceExvat:   assetData["price_exvat"].String(),
		Downloaded:   assetData["downloaded"].String(),
		Date:         convertedDate,
		Quantity:     assetData["quantity"].String(),
		Reason:       assetData["reason"].String(),
		OtherLicense: assetData["other_license"].String(),
		Package:      assetData["package"].String(),
		Currency:     assetData["currency"].String(),
		Refunded:     assetData["refunded"].String(),
		Invoice:      assetData["invoice"].String(),
	}

	maskedInvoice := MaskLeft(func() string { var str = asset.Invoice; return string(str) }())
	Log.Infof("Successfully verified Asset: %s: for invoice %s", asset.Package, maskedInvoice)
	//return asset
}
