package dns

import (
	"net"
)

// WhoisResult contains nameserver data.
type WhoisResult struct {
	NameServers []string `json:"name_servers"`
	Error       string   `json:"error,omitempty"`
}

// GetWhois retrieves NS records.
func GetWhois(domain string) WhoisResult {
	res := WhoisResult{}
	
	result, err := net.LookupNS(domain)
	if err != nil {
		res.Error = err.Error()
		return res
	}
	
	for _, ns := range result {
		res.NameServers = append(res.NameServers, ns.Host)
	}
	return res
}
