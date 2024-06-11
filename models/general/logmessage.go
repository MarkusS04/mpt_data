// Package generalmodel provides types for db and api
package generalmodel

// Log message for Info
const (
	StartExecPDFAutoremoval = "Start execution of pdf autoremoval"
	EndExecPDFAutoremoval   = "End execution of pdf autoremoval"

	DBMigrated = "database migration succesfull"

	UnkownAcceptHeader = "unknown accept header"
)

// Log message for Warning
const (
	UserInvalidLogin = "Invalid credentials"
)

// Log message for Error
const (
	APIStartFailed = "failed to start API"

	DBMigrationFailed  = "database migration failed"
	DBLoadDataFailed   = "failed to load data from db"
	DBSaveDataFailed   = "failed to store data in db"
	DBUpdateDataFailed = "failed to upate data in db"
	DBDeleteDataFailed = "failed to delete data from db"

	UserCreationFailed = "could not create user"

	PDFRemovalFailed      = "could not delete pdf-file"
	PDFFileCreationFailed = "could not store pdf-file"

	PlanCreationFailed = "failed to create plan"
	PlanCreationError  = "error during plan creation"

	InternalError    = "Internal Server Error"
	StatusBadRequest = "Bad Request"
)

// Keys for logger
const (
	AdditionalInfo = "AdditionalInfo"
	AcceptHeader   = "AcceptHedear"
)
