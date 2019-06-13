package models

import (
	"fmt"
	"time"

	"github.com/parnurzeal/gorequest"
	"github.com/tidwall/gjson"

	. "github.com/instance-id/GoVerifier-dgo/utils"
)

var log = Log

type InvoiceDataAccessObject struct{}

type UserPackagesDetail struct {
	VerifiedUser `xorm:"extends"`
	UserPackages `xorm:"extends"`
}

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

var InvoiceDAO *InvoiceDataAccessObject

func (i *InvoiceDataAccessObject) CheckInvoice(invoice string) bool {
	exists, err := Dba.Table(UserPackagesDAO.TableName()).Where("invoice = ?", invoice).Exist()
	LogInfof("Database connection failed to return initial username check ", err)
	return exists
}

// --- Using asset store API, verify that the provided invoice was a valid purchase -----------------------------------------------------------------
func (i *InvoiceDataAccessObject) VerifyInvoice(invoiceNum string, assetChoice string) (bool, string, *Invoices) {
	assetName := Dac.Assets.Packages[assetChoice]
	apiKey := Dac.Assets.ApiKeys[assetChoice]

	url := "https://api.assetstore.unity3d.com:443/publisher/v1/invoice/verify.json"

	request := gorequest.New()
	resp, body, err := request.Set("Accept", "application/json").Get(url).Query("key=" + apiKey + "&" + "invoice=" + invoiceNum).End()
	if err != nil {
		msg := fmt.Sprintf("Received warning from Server: %s", err)
		Log.Errorf(msg, err)
		return false, msg, nil
	}
	if len(body) <= 0 {
		msg := fmt.Sprintf("Received empty reply from server. Status code received: %s", resp.Status)
		Log.Errorf(msg, resp.StatusCode)
		return false, msg, nil
	}

	assetData := gjson.Get(body, "invoices.#[package=="+assetName+"]").Map()
	Log.Infof("Data %v", assetData)

	if len(assetData) <= 1 {
		msg := fmt.Sprintf("Could not verify: %s for invoice %s. Please check the invoice number and try again.", assetName, invoiceNum)
		Log.Errorf(msg, nil)
		return false, msg, nil
	}

	convertedDate := ConvertTime(assetData["date"].String())
	trimmedInvoice := TrimInvoice(assetData["invoice"].String())

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
		Invoice:      trimmedInvoice,
	}

	maskedInvoice := MaskLeft(func() string { var str = asset.Invoice; return string(str) }())
	msg := fmt.Sprintf("Successfully verified Asset: %s: for invoice %s", asset.Package, maskedInvoice)
	Log.Debugf(msg)
	return true, msg, &asset
}

func (i *InvoiceDataAccessObject) AddInvoice(user *VerifiedUser, packages *UserPackages) (bool, string) {
	Log.Debugf("Username: %s", user.Username) // ---------------------------------------------------------------------------------------- TODO Remove

	exists, err := Dba.Table(VerifiedUserDAO.TableName()).Where("username = ?", user.Username).Exist()
	hadError := LogErrorRet("Database connection failed to return initial username check ", err)
	if hadError {
		return true, fmt.Sprintf("Database connection failed to return initial username check %s", err)
	}

	Log.Debugf("Exists?: %t", exists) // ------------------------------------------------------------------------------------------------ TODO Remove

	if exists {
		_, err := Dba.Table(VerifiedUserDAO.TableName()).Where("username = ?", user.Username).Cols("id").Get(&packages.Id)
		hadError := LogErrorRet("Unable to query user id: ", err)
		if hadError {
			return true, fmt.Sprintf("Database connection failed to return ID: %s ", err)
		}

		Log.Debugf("Package Id: %v", packages.Id) // ------------------------------------------------------------------------------------ TODO Remove

		_, err = Dba.Table(UserPackagesDAO.TableName()).InsertOne(packages)
		hadError = LogErrorRet("Unable to add packages to user: ", err)
		if hadError {
			return true, fmt.Sprintf("Database connection failed to return package ID: %s ", err)
		}

		return false, fmt.Sprintf("Package: %s has been successfully applied to %s ", packages.Packages, user.Username)
	}

	Log.Debugf("Begin adding to database - Id: %v", user.Username) // ------------------------------------------------------------------- TODO Remove

	// --- Add new user to database ----------------------------------------------------------------------------------------
	_, err = Dba.Table(VerifiedUserDAO.TableName()).InsertOne(user)
	hadError = LogErrorRet("Unable to insert user: ", err)
	if hadError {
		return true, fmt.Sprintf("System encountered communication error. Please contact support:  %s", err)
	}

	Log.Debugf("User added to database - Id: %v", user.Id) // --------------------------------------------------------------------------- TODO Remove

	// --- Add verified package to new user --------------------------------------------------------------------------------
	_, err = Dba.Table(UserPackagesDAO.TableName()).InsertOne(func() *UserPackages { packages.Id = user.Id; return packages }())
	hadError = LogErrorRet("Unable to insert package: ", err)
	if hadError {
		return true, fmt.Sprintf("System encountered communication error. Please contact support:  %s", err)
	}

	Log.Debugf("Package added to database - Id: %v", packages.Id) // -------------------------------------------------------------------- TODO Remove
	return false, fmt.Sprintf("User successfully added: UserId: %v  Packages: %v ", user.Username, packages.Packages)
}
