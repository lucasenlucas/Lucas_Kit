package dns

import (
	"context"
	"strings"
	"time"

	"github.com/miekg/dns"
)

// EmailSecurityResult contains DKIM, SPF, DMARC, etc info.
type EmailSecurityResult struct {
	SPF    string   `json:"spf,omitempty"`
	DMARC  string   `json:"dmarc,omitempty"`
	DKIM   []string `json:"dkim,omitempty"`
	MX     []string `json:"mx,omitempty"`
	TLSRPT string   `json:"tls_rpt,omitempty"`
	MTASTS string   `json:"mta_sts,omitempty"`
	Error  string   `json:"error,omitempty"`
}

// GetEmailSecurity fetches mail-related DNS components.
func GetEmailSecurity(ctx context.Context, domain string, timeout time.Duration) EmailSecurityResult {
	resolver := pickResolver("")
	client := new(dns.Client)
	client.Timeout = timeout

	res := EmailSecurityResult{}

	// SPF
	txt, err := queryType(ctx, client, resolver, domain, dns.TypeTXT)
	if err == nil {
		res.SPF = findTXTContains(txt, "v=spf1")
	}

	// DMARC
	dmarc, err := queryType(ctx, client, resolver, "_dmarc."+domain, dns.TypeTXT)
	if err == nil {
		res.DMARC = findTXTContains(dmarc, "v=DMARC1")
	}

	// DKIM
	dkimSelectors := []string{"default", "selector1", "selector2", "s1", "s2", "k1", "google"}
	for _, sel := range dkimSelectors {
		d, err := queryType(ctx, client, resolver, sel+"._domainkey."+domain, dns.TypeTXT)
		if err == nil {
			v := findTXTContains(d, "v=DKIM1")
			if v != "" {
				res.DKIM = append(res.DKIM, sel+": "+v)
			}
		}
	}

	// MX
	mx, err := queryType(ctx, client, resolver, domain, dns.TypeMX)
	if err == nil {
		for _, rr := range mx {
			if m, ok := rr.(*dns.MX); ok {
				res.MX = append(res.MX, strings.TrimSuffix(m.Mx, "."))
			}
		}
	}

	// TLS-RPT
	tlsRpt, err := queryType(ctx, client, resolver, "_smtp._tls."+domain, dns.TypeTXT)
	if err == nil {
		res.TLSRPT = findTXTContains(tlsRpt, "v=TLSRPTv1")
	}

	// MTA-STS
	mtaSts, err := queryType(ctx, client, resolver, "_mta-sts."+domain, dns.TypeTXT)
	if err == nil {
		res.MTASTS = findTXTContains(mtaSts, "v=STSv1")
	}

	return res
}

func findTXTContains(rrs []dns.RR, needle string) string {
	needle = strings.ToLower(needle)
	for _, rr := range rrs {
		t, ok := rr.(*dns.TXT)
		if !ok {
			continue
		}
		joined := strings.Join(t.Txt, "")
		if strings.Contains(strings.ToLower(joined), needle) {
			return joined
		}
	}
	return ""
}
