package utils

import (
	"net"
	"net/http"
	"time"
)

// DoctorResult holds the findings of environment checks.
type DoctorResult struct {
	Internet  bool
	DNS       bool
	Resolvers []string
	Issues    []string
}

// RunDoctor performs self-checks to ensure NetScope can operate fully.
func RunDoctor() DoctorResult {
	res := DoctorResult{
		Resolvers: []string{},
		Issues:    []string{},
	}

	// Check HTTP connectivity
	client := http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get("http://1.1.1.1")
	if err == nil {
		res.Internet = true
		resp.Body.Close()
	} else {
		res.Issues = append(res.Issues, "No basic HTTP internet connectivity.")
	}

	// Check DNS
	addrs, err := net.LookupHost("google.com")
	if err == nil && len(addrs) > 0 {
		res.DNS = true
	} else {
		res.Issues = append(res.Issues, "DNS resolution failed.")
	}

	return res
}
