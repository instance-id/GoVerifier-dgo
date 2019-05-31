package components

type DB interface {
	// --- Quick lookup if invoice exists -------------------------------------------------------------------
	FindInvoice(invoice string) (string, error)
	// --- Add new invoice ----------------------------------------------------------------------------------
	AddInvoice(username string, invoice string, pkg string, purdate string, email string) (string, error)
	// --- Remove invoice from user -------------------------------------------------------------------------
	DeleteInvoice(invoice string) (string, error)
	// --- Search and return invoice details ----------------------------------------------------------------
	SearchInvoice(invoice string) (string, error)
	// --- Search and return user details -------------------------------------------------------------------
	SearchUser(username string) (string, error)
	// --- Initial DB setup ---------------------------------------------------------------------------------
	DbSetup(setup bool) (bool, string, error)
	// --- Ensure that all migrations are done --------------------------------------------------------------
	Ensure() error
	// --- Run db mainloop ----------------------------------------------------------------------------------
	Run()
	// --- Close access to database -------------------------------------------------------------------------
	Close() error
}
