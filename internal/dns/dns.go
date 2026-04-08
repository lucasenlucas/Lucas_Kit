package dns

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/miekg/dns"
)

// DNSResult contains requested DNS records.
type DNSResult struct {
	Records map[string][]string `json:"records"`
	Error   string              `json:"error,omitempty"`
}

// GetDNS fetches A, AAAA, MX, NS, and TXT records.
func GetDNS(ctx context.Context, domain string, timeout time.Duration) DNSResult {
	resolver := pickResolver("")
	client := new(dns.Client)
	client.Timeout = timeout

	res := DNSResult{
		Records: make(map[string][]string),
	}

	queries := []struct {
		name  string
		qtype uint16
	}{
		{"A", dns.TypeA},
		{"AAAA", dns.TypeAAAA},
		{"MX", dns.TypeMX},
		{"NS", dns.TypeNS},
		{"TXT", dns.TypeTXT},
	}

	for _, qu := range queries {
		rrs, err := queryType(ctx, client, resolver, domain, qu.qtype)
		if err != nil {
			continue // skip errors to try other records
		}
		
		var recStrings []string
		for _, rr := range rrs {
			if rr.Header() != nil && rr.Header().Rrtype == dns.TypeOPT {
				continue
			}
			recStrings = append(recStrings, rr.String())
		}
		if len(recStrings) > 0 {
			res.Records[qu.name] = recStrings
		}
	}

	return res
}

func pickResolver(flagVal string) string {
	if flagVal != "" {
		return flagVal
	}
	if cfg, err := dns.ClientConfigFromFile("/etc/resolv.conf"); err == nil && len(cfg.Servers) > 0 {
		return net.JoinHostPort(cfg.Servers[0], cfg.Port)
	}
	return "8.8.8.8:53"
}

func queryType(ctx context.Context, client *dns.Client, resolver, name string, qtype uint16) ([]dns.RR, error) {
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(name), qtype)
	m.RecursionDesired = true

	rctx, cancel := context.WithTimeout(ctx, client.Timeout)
	defer cancel()

	in, _, err := client.ExchangeContext(rctx, m, resolver)
	if err != nil {
		return nil, err
	}
	if in.Rcode != dns.RcodeSuccess && in.Rcode != dns.RcodeNameError {
		return nil, fmt.Errorf("dns rcode %s", dns.RcodeToString[in.Rcode])
	}
	
	var out []dns.RR
	out = append(out, in.Answer...)
	out = append(out, in.Extra...)
	return out, nil
}
