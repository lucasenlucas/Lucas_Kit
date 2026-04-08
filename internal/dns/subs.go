package dns

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// SubsResult contains subdomains found via crt.sh
type SubsResult struct {
	Subdomains []string `json:"subdomains"`
	Error      string   `json:"error,omitempty"`
}

// GetSubdomains queries crt.sh for known subdomains.
func GetSubdomains(domain string) SubsResult {
	res := SubsResult{}
	
	client := http.Client{Timeout: 15 * time.Second}
	url := fmt.Sprintf("https://crt.sh/?q=%%.%s&output=json", domain)
	
	resp, err := client.Get(url)
	if err != nil {
		res.Error = err.Error()
		return res
	}
	defer resp.Body.Close()

	var results []struct {
		NameValue string `json:"name_value"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		res.Error = "Failed to parse crt.sh response"
		return res
	}

	subsMap := make(map[string]bool)
	for _, crtRec := range results {
		names := strings.Split(crtRec.NameValue, "\n")
		for _, name := range names {
			name = strings.TrimSpace(name)
			if name != "" && !strings.Contains(name, "*") {
				subsMap[name] = true
			}
		}
	}

	for sub := range subsMap {
		res.Subdomains = append(res.Subdomains, sub)
	}
	
	return res
}
