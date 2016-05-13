package types

import (
    "sort"
)

// Service contains service definition
type Service struct {
    ID []string
    ProxyServices []string
}

// HasProxyService checks if given serviceKey exists
func (s Service) HasProxyService(serviceKey string) bool {
    i := sort.SearchStrings(s.ProxyServices, serviceKey)
    
    return i < len(s.ProxyServices) && s.ProxyServices[i] == serviceKey
}