package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

func runWebAnalysis(o options) {
	domain := normalizeDomain(o.domain)

	client := &http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	var doReq = func() (*http.Response, []*http.Response, error) {
		clientWithRedirects := &http.Client{Timeout: 10 * time.Second}
		req, _ := http.NewRequest("GET", "https://"+domain, nil)
		var redirects []*http.Response
		clientWithRedirects.CheckRedirect = func(r *http.Request, via []*http.Request) error {
			redirects = append(redirects, r.Response)
			if len(via) >= 10 {
				return fmt.Errorf("stopped after 10 redirects")
			}
			return nil
		}

		req.Header.Set("User-Agent", "NetScope-Analysis/1.0")
		resp, err := clientWithRedirects.Do(req)
		if err != nil {
			req, _ = http.NewRequest("GET", "http://"+domain, nil)
			resp, err = clientWithRedirects.Do(req)
		}
		return resp, redirects, err
	}

	fn := strings.ToLower(o.function)

	switch fn {
	case "status":
		resp, _, err := doReq()
		if err != nil {
			fmt.Printf("[!] HTTP Error: %v\n", err)
			return
		}
		if !o.jsonOut {
			fmt.Printf("│ [STATUS] Code: %d\n", resp.StatusCode)
		}
	case "redirect":
		_, redirects, err := doReq()
		if err != nil {
			fmt.Printf("[!] HTTP Error: %v\n", err)
			return
		}
		if !o.jsonOut {
			fmt.Println("│ [REDIRECT] Keten:")
			for i, r := range redirects {
				fmt.Printf("│   %d. %s -> %s\n", i+1, r.Request.URL, r.Header.Get("Location"))
			}
		}
	case "ssl":
		resp, _, err := doReq()
		if err != nil {
			fmt.Printf("[!] HTTP Error: %v\n", err)
			return
		}
		if resp != nil && resp.TLS != nil && len(resp.TLS.PeerCertificates) > 0 {
			cert := resp.TLS.PeerCertificates[0]
			if !o.jsonOut {
				fmt.Printf("│ [SSL] Issuer:     %s\n", cert.Issuer.CommonName)
				fmt.Printf("│ [SSL] Geldig tot: %v\n", cert.NotAfter.Format("02-01-2006"))
				fmt.Printf("│ [SSL] Domeinen:   %v\n", cert.DNSNames)
			}
		} else {
			fmt.Println("[!] Geen SSL/TLS certificaat gevonden.")
		}
	case "headers":
		resp, _, err := doReq()
		if err != nil {
			fmt.Printf("[!] HTTP Error: %v\n", err)
			return
		}
		if !o.jsonOut {
			headers := []string{"Strict-Transport-Security", "Content-Security-Policy", "X-Frame-Options", "X-Content-Type-Options", "Referrer-Policy"}
			fmt.Println("│ [HEADERS] Security Headers:")
			for _, h := range headers {
				val := resp.Header.Get(h)
				if val == "" {
					val = "❌ ONTBREKT"
				} else {
					val = "✅ " + val
				}
				fmt.Printf("│   %-25s: %s\n", h, val)
			}
		}
	case "robots":
		resp, err := client.Get("https://" + domain + "/robots.txt")
		if err != nil {
			resp, err = client.Get("http://" + domain + "/robots.txt")
		}
		if err == nil && resp.StatusCode == 200 {
			body, _ := io.ReadAll(resp.Body)
			if !o.jsonOut {
				fmt.Println("│ [ROBOTS] Bestandsinhoud:")
				lines := strings.Split(string(body), "\n")
				for _, line := range lines {
					if strings.TrimSpace(line) != "" {
						fmt.Printf("│   %s\n", line)
					}
				}
			}
			resp.Body.Close()
		} else {
			fmt.Println("[!] Geen robots.txt gevonden.")
		}
	case "sitemap":
		resp, err := client.Get("https://" + domain + "/sitemap.xml")
		if err != nil {
			resp, err = client.Get("http://" + domain + "/sitemap.xml")
		}
		if err == nil && resp.StatusCode == 200 {
			if !o.jsonOut {
				fmt.Println("│ [SITEMAP] Sitemap gevonden op /sitemap.xml")
			}
			resp.Body.Close()
		} else {
			fmt.Println("[!] Geen sitemap.xml gevonden.")
		}
	case "tech":
		resp, _, err := doReq()
		if err == nil {
			techList := []string{}
			serverHeader := resp.Header.Get("Server")
			if serverHeader != "" {
				techList = append(techList, serverHeader)
			}
			bodyBytes, _ := io.ReadAll(resp.Body)
			bodyStr := strings.ToLower(string(bodyBytes))
			if strings.Contains(bodyStr, "wp-content") {
				techList = append(techList, "WordPress")
			}
			if strings.Contains(bodyStr, "react") {
				techList = append(techList, "React")
			}
			if !o.jsonOut {
				fmt.Printf("│ [TECH] Gedetecteerd: %s\n", strings.Join(techList, ", "))
			}
		}
	case "forms":
		resp, _, err := doReq()
		if err == nil {
			forms := getPageForms(resp.Body)
			if !o.jsonOut {
				fmt.Println("│ [FORMS] Gevonden formulieren:")
				for _, f := range forms {
					fmt.Printf("│   - %s\n", f)
				}
			}
		}
	case "links":
		resp, _, err := doReq()
		if err == nil {
			links := getPageLinks(resp.Body)
			if !o.jsonOut {
				fmt.Printf("│ [LINKS] %d links gevonden:\n", len(links))
				for _, l := range links {
					fmt.Printf("│   - %s\n", l)
				}
			}
		}
	case "scripts":
		resp, _, err := doReq()
		if err == nil {
			scripts := getPageScripts(resp.Body)
			if !o.jsonOut {
				fmt.Printf("│ [SCRIPTS] %d externe scripts gevonden:\n", len(scripts))
				for _, s := range scripts {
					fmt.Printf("│   - %s\n", s)
				}
			}
		}
	case "cookies":
		resp, _, err := doReq()
		if err == nil {
			if !o.jsonOut {
				fmt.Println("│ [COOKIES] Analyse:")
				for _, c := range resp.Cookies() {
					fmt.Printf("│   - %-15s | Secure: %-5v | HttpOnly: %-5v | SameSite: %v\n", c.Name, c.Secure, c.HttpOnly, c.SameSite)
				}
			}
		}
	case "ports":
		commonPorts := []int{80, 443, 8080, 8443}
		if !o.jsonOut {
			fmt.Println("│ [PORTS] Scan:")
			for _, p := range commonPorts {
				addr := net.JoinHostPort(domain, fmt.Sprintf("%d", p))
				conn, err := net.DialTimeout("tcp", addr, 1*time.Second)
				if err == nil {
					fmt.Printf("│   - Port %d: ✅ OPEN\n", p)
					conn.Close()
				} else {
					fmt.Printf("│   - Port %d: ❌ GESLOTEN\n", p)
				}
			}
		}
	case "favicon":
		resp, _, err := doReq()
		if err == nil {
			fav := getFaviconURL(domain, resp.Body)
			if !o.jsonOut {
				fmt.Printf("│ [FAVICON] URL: %s\n", fav)
			}
		}
	case "title":
		resp, _, err := doReq()
		if err == nil {
			title := getPageTitle(resp.Body)
			if !o.jsonOut {
				fmt.Printf("│ [TITLE] Paginatitel: %s\n", title)
			}
		}
	}
}
