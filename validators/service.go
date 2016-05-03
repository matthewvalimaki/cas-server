package validators

import (   
    "fmt"
    "errors"
    "sort"
    
    "github.com/matthewvalimaki/cas-server/types"
)

// ValidateService validates service
func ValidateService(ID string, config *types.Config) (*types.CasError) {
    err := validateServiceIDLength(ID)    
    if (err != nil) {
        return err
    }
    
    err = validateServiceID(ID, config)    
    if (err != nil) {
        return err
    }

    return nil
}

func validateServiceIDLength(ID string) *types.CasError {
    if len(ID) == 0 {
        return &types.CasError{Error: errors.New("Required query parameter `service` was not defined."), CasErrorCode: types.CAS_ERROR_CODE_INVALID_SERVICE}
    }  
    
    return nil
}

func validateServiceID(ID string, config *types.Config) *types.CasError {
    sort.Strings(config.FlatServiceIDList)
    
    i := sort.SearchStrings(config.FlatServiceIDList, ID)
    if i < len(config.FlatServiceIDList) && config.FlatServiceIDList[i] == ID {
        return nil
    }

    return &types.CasError{Error: fmt.Errorf("Unknown service `%s`.", ID), CasErrorCode: types.CAS_ERROR_CODE_INVALID_SERVICE}
}