package apperror

// General codes
const (
	CodeInternalError   = "INTERNAL_ERROR"
	CodeInvalidInput    = "INVALID_INPUT"
	CodeUnauthorized    = "UNAUTHORIZED"
	CodeForbidden       = "FORBIDDEN"
	CodeNotImplemented  = "NOT_IMPLEMENTED"
	CodeTimeout         = "TIMEOUT"
	CodeDependencyError = "DEPENDENCY_ERROR"
)

// Resource-specific codes
const (
	CodeUserNotFound     = "USER_NOT_FOUND"
	CodeUserConflict     = "USER_CONFLICT"
	CodeResourceNotFound = "RESOURCE_NOT_FOUND"
	CodeResourceConflict = "RESOURCE_CONFLICT"
)

// DB-related
const (
	CodeDBError      = "DB_ERROR"
	CodeDBNoRows     = "DB_NO_ROWS"
	CodeDBConnFailed = "DB_CONN_FAILED"
	CodeDBTxFailed   = "DB_TX_FAILED"
	CodeDBConstraint = "DB_CONSTRAINT"
)

// Tambahan
const (
	CodeRoleNotFound = "ROLE_NOT_FOUND"
	CodeAuthNotFound = "AUTH_NOTFOUND"
)
