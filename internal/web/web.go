package web

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"netscope/internal/crawl" // Ensure this matches module path later
)

type WebResult struct {
	Status   int                 `json:"status_code,omitempty"`
	Redirect []string            `json:"redirect_chain,omitempty"`
	TLS      map[string]string   `json:"tls,omitempty"`
	Headers  map[string]string   `json:"security_headers,omitempty"`
	Cookies  []map[string]string `json:"cookies,omitempty"`
	Tech     []string            `json:"technology,omitempty"`
	Ports    map[int]bool        `json:"ports_open,omitempty"`
	Robots   string              `json:"robots_txt,omitempty"`
	Sitemap  bool                `json:"sitemap_found,omitempty"`
	Title    string              `json:"title,omitempty"`
	Favicon  string              `json:"favicon,omitempty"`
	Crawl    *crawl.CrawlResult  `json:"crawl_data,omitempty"`
	Error    string              `json:"error,omitempty"`
}

func AnalyzeWeb(domain string, modules []string) WebResult {
	res := WebResult{}

	// Helper to make a standard request that tracks redirects
	clientWithRedirects := &http.Client{Timeout: 10 * time.Second}
	var redirects []*http.Response
	clientWithRedirects.CheckRedirect = func(r *http.Request, via []*http.Request) error {
		redirects = append(redirects, r.Response)
		if len(via) >= 10 {
			return fmt.Errorf("stopped after 10 redirects")
		}
		return nil
	}

	req, _ := http.NewRequest("GET", "https://"+domain, nil)
	req.Header.Set("User-Agent", "NetScope-Analysis/2.0")
	resp, err := clientWithRedirects.Do(req)

	if err != nil {
		req, _ = http.NewRequest("GET", "http://"+domain, nil)
		resp, err = clientWithRedirects.Do(req)
	}

	if err != nil {
		res.Error = fmt.Sprintf("HTTP Error: %v", err)
		return res
	}
	defer resp.Body.Close()

	// Cache body for crawl tools
	bodyBytes, _ := io.ReadAll(resp.Body)

	for _, mod := range modules {
		switch mod {
		case "status":
			res.Status = resp.StatusCode
		case "redirect":
			for i, r := range redirects {
				res.Redirect = append(res.Redirect, fmt.Sprintf("%d. %s -> %s", i+1, r.Request.URL, r.Header.Get("Location")))
			}
		case "tls", "ssl":
			if resp.TLS != nil && len(resp.TLS.PeerCertificates) > 0 {
				cert := resp.TLS.PeerCertificates[0]
				res.TLS = map[string]string{
					"Issuer":      cert.Issuer.CommonName,
					"Valid_Until": cert.NotAfter.Format("2006-01-02"),
					"DNS_Names":   strings.Join(cert.DNSNames, ", "),
				}
			}
		case "headers":
			res.Headers = make(map[string]string)
			headersToCheck := []string{"Strict-Transport-Security", "Content-Security-Policy", "X-Frame-Options", "X-Content-Type-Options", "Referrer-Policy"}
			for _, h := range headersToCheck {
				val := resp.Header.Get(h)
				if val == "" {
					res.Headers[h] = "MISSING"
				} else {
					res.Headers[h] = val
				}
			}
		case "cookies":
			for _, c := range resp.Cookies() {
				res.Cookies = append(res.Cookies, map[string]string{
					"Name":      c.Name,
					"Secure":    fmt.Sprintf("%v", c.Secure),
					"HttpOnly":  fmt.Sprintf("%v", c.HttpOnly),
					"SameSite":  fmt.Sprintf("%v", c.SameSite),
				})
			}
		case "tech":
			res.Tech = analyzeTech(resp, bodyBytes)
		case "ports":
			res.Ports = scanPorts(domain)
		case "robots":
			robCode, robBody := fetchURL(fmt.Sprintf("https://%s/robots.txt", domain))
			if robCode != 200 {
				robCode, robBody = fetchURL(fmt.Sprintf("http://%s/robots.txt", domain))
			}
			if robCode == 200 {
				res.Robots = robBody
			}
		case "sitemap":
			smCode, _ := fetchURL(fmt.Sprintf("https://%s/sitemap.xml", domain))
			if smCode != 200 {
				smCode, _ = fetchURL(fmt.Sprintf("http://%s/sitemap.xml", domain))
			}
			res.Sitemap = smCode == 200
		case "title":
			res.Title = crawl.GetTitle(bytes.NewReader(bodyBytes))
		case "favicon":
			res.Favicon = crawl.GetFavicon(domain, bytes.NewReader(bodyBytes))
		case "links":
			ensureCrawlData(&res)
			res.Crawl.Links = crawl.GetLinks(bytes.NewReader(bodyBytes))
		case "scripts":
			ensureCrawlData(&res)
			res.Crawl.Scripts = crawl.GetScripts(bytes.NewReader(bodyBytes))
		case "forms":
			ensureCrawlData(&res)
			res.Crawl.Forms = crawl.GetForms(bytes.NewReader(bodyBytes))
		}
	}

	return res
}

func ensureCrawlData(res *WebResult) {
	if res.Crawl == nil {
		res.Crawl = &crawl.CrawlResult{}
	}
}

func analyzeTech(resp *http.Response, body []byte) []string {
	var tech []string
	if srv := resp.Header.Get("Server"); srv != "" {
		tech = append(tech, srv)
	}
	bodyStr := strings.ToLower(string(body))
	if strings.Contains(bodyStr, "wp-content") {
		tech = append(tech, "WordPress")
	}
	if strings.Contains(bodyStr, "react") {
		tech = append(tech, "React")
	}
	return tech
}

func scanPorts(domain string) map[int]bool {
	res := make(map[int]bool)
	ports := []int{80, 443, 8080, 8443}
	for _, p := range ports {
		addr := fmt.Sprintf("%s:%d", domain, p)
		conn, err := net.DialTimeout("tcp", addr, 1*time.Second)
		if err == nil {
			res[p] = true
			conn.Close()
		} else {
			res[p] = false
		}
	}
	return res
}

func fetchURL(url string) (int, string) {
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return 0, ""
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, string(b)
}
