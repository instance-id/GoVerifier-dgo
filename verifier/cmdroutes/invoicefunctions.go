// --- cmd specific functions which are best to remain within the package ---------------------------------------------------------------------------
package cmdroutes

import (
	"fmt"

	. "github.com/instance-id/GoVerifier-dgo/data"
	. "github.com/instance-id/GoVerifier-dgo/utils"
)

// --- Generic input text query -------------------------------------------------------------------
func GenericQuery(l LocalContext, prompt string) string {
	return func(string) string {
		query, err := QueryInput(l.C, l.Ctx, prompt, TimeoutDuration)
		result := InputWarn("Input not received: ", err)
		if result == true {
			return "Input not received: Action timed out"
		}
		return query.Content
	}(prompt)
}

// --- Queries user for invoice input ------------------------------------------------------------- TODO Test Function
func GetInvoice(l LocalContext, prompt string) string {
	return func(string) string {
		invoice, err := QueryInput(l.C, l.Ctx, prompt, TimeoutDuration)
		result := InputWarn("Input not received: ", err)
		if result == true {
			return "Input not received: Action timed out"
		}
		return TrimInvoice(invoice.Content)
	}(prompt)
}

// --- Obtain asset identification ----------------------------------------------------------------
func EnterAsset(l LocalContext, prompt string) string {
	return func(string) string {
		asset, err := QueryInput(l.C, l.Ctx, prompt, TimeoutDuration)
		result := InputWarn("Input not received: ", err)
		if result == true {
			return "Input not received: Action timed out"
		}
		return asset.Content
	}(prompt)
}

// --- Obtain invoice number ----------------------------------------------------------------------
func EnterInvoice(l LocalContext, asset string, assetDesc string) string {
	prompt := fmt.Sprintf("Please enter invoice number for %s: %s", assetDesc, asset)

	return func(string) string {
		invoice, err := QueryInput(l.C, l.Ctx, prompt, TimeoutDuration)
		result := InputWarn("Input not received: ", err)
		if result == true {
			return "Input not received: Action timed out"
		}
		return TrimInvoice(invoice.Content)
	}(prompt)
}

// --- Obtain email address ----------------------------------------------------------------------
func EnterEmail(l LocalContext) string {
	prompt := fmt.Sprintf("Please enter email address associated with asset purchase")

	return func(string) string {
		email, err := QueryInput(l.C, l.Ctx, prompt, TimeoutDuration)
		result := InputWarn("Input not received: ", err)
		if result == true {
			return "Input not received: Action timed out"
		}
		return email.Content
	}(prompt)
}
