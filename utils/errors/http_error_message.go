package errors

// custom error message
const (
	JSONParseError = "JSON Parse Error"
	// Create account
	NotValidEmail        = "not valid email"
	PasswordError        = "password error"
	InternalError        = "internal error"
	InsertDBError        = "insert DB error"
	AccountAlreadyExists = "account already exists"

	// Login error message
	AccountNotExist        = "account not exist"
	PasswordDecryptedError = "password decrypt error"
	PasswordIncorrect      = "password incorrect"

	// rcic error message
	CreateHTTPRequestFailed   = "create HTTP request failed"
	RequestRemoteServerFailed = "request remote server failed"
	ReadResponseBodyFailed    = "read response body failed"
)
