package types

// CasErrorCode type declaration
type CasErrorCode int

const (
    // CAS_ERROR_CODE_INVALID_REQUEST "not all of the required request parameters were present"
    CAS_ERROR_CODE_INVALID_REQUEST CasErrorCode = 1 + iota
    // CAS_ERROR_CODE_INVALID_TICKET_SPEC "failure to meet the requirements of validation specification"
    CAS_ERROR_CODE_INVALID_TICKET_SPEC
    // CAS_ERROR_CODE_INVALID_TICKET "the ticket provided was not valid, or the ticket did not come from an initial login and renew was set on validation."
    CAS_ERROR_CODE_INVALID_TICKET
    // INVALID_SERVICE "the ticket provided was valid, but the service specified did not match the service associated with the ticket."
    CAS_ERROR_CODE_INVALID_SERVICE
    // CAS_ERROR_CODE_INTERNAL_ERROR "an internal error occurred during ticket validation"
    CAS_ERROR_CODE_INTERNAL_ERROR
)

// CasErrorCodes contains all error codes in string format
var CasErrorCodes = [...]string {
    "INVALID_REQUEST",
    "INVALID_TICKET_SPEC",
    "INVALID_TICKET",
    "INVALID_SERVICE",
    "INTERNAL_ERROR",
}

func (casErrorCode CasErrorCode) String() string {
 return CasErrorCodes[casErrorCode - 1]
}

// CasError contains CAS error information
type CasError struct {
    Error error
    CasErrorCode CasErrorCode
}