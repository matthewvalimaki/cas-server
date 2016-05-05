package validators

import (   
    "errors"
    
    "github.com/matthewvalimaki/cas-server/types"
)

// ValidateService validates service
func ValidateService(ID string, config *types.Config) (*types.CasError) {
    err := validateServiceIDLength(ID)    
    if err != nil {
        return err
    }
    
    err = validateServiceID(ID, config)    
    if err != nil {
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

func validateServiceID(serviceID string, config *types.Config) *types.CasError {
    if _, ok := config.FlatServiceIDList[serviceID]; ok {
        return nil
    }

    return &types.CasError{Error: errors.New("Unknown service."), CasErrorCode: types.CAS_ERROR_CODE_INVALID_SERVICE}
}