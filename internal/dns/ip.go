package dns

import (
	"net"
)

// IPResult contains the IPs resolved for a domain.
type IPResult struct {
	IPs   []string `json:"ips"`
	Error string   `json:"error,omitempty"`
}

// GetIP resolves a domain to its IP addresses.
func GetIP(domain string) IPResult {
	res := IPResult{}
	
	ips, err := net.LookupIP(domain)
	if err != nil {
		res.Error = err.Error()
		return res
	}
	
	for _, ip := range ips {
		res.IPs = append(res.IPs, ip.String())
	}
	return res
}
