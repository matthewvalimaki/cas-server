package validators

import (
    "errors"
    
    "github.com/matthewvalimaki/cas-server/types"
)

// ValidateFormat validates that format is of correct value
func ValidateFormat(format string) *types.CasError {
    if format != "XML" && format != "JSON" {
        return &types.CasError{Error: errors.New("Query parameter `format` contained illegal value. Allowed values are `XML' and `JSON`."), CasErrorCode: types.CAS_ERROR_CODE_INVALID_PROXY_CALLBACK}
    }
    
    return nil
}