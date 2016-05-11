package validators

import ( 
    "fmt" 
    "errors"
    "regexp"
    
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
    for supportedServiceID := range config.FlatServiceIDList {
        if matched, _ := regexp.MatchString(supportedServiceID, serviceID); matched {
            return nil
        }
    }

    return &types.CasError{Error: fmt.Errorf("Unknown service `%s`", serviceID), CasErrorCode: types.CAS_ERROR_CODE_INVALID_SERVICE}
}