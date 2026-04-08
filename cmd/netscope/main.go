package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"netscope/internal/dns"
	"netscope/internal/output"
	"netscope/internal/utils"
	"netscope/internal/web"
)

const version = "8.4.2026"

var profileMap = map[string][]string{
	"quick":      {"ip", "status", "headers", "tls", "title"},
	"web":        {"status", "headers", "redirect", "cookies", "tech", "forms", "scripts", "links"},
	"dns-full":   {"dns", "whois", "ip", "email", "subs"},
	"full":       {"dns", "whois", "ip", "email", "subs", "status", "headers", "redirect", "cookies", "tech", "forms", "scripts", "links", "robots", "sitemap", "title", "favicon", "ports"},
	"pageinfo":   {"title", "favicon"},
	"web-basic":  {"status", "headers", "redirect"},
	"crawl-lite": {"links", "scripts", "forms"},
}

// Module map for CLI help
var allModules = []string{
	"dns", "whois", "ip", "subs", "email", "ports",
	"status", "headers", "redirect", "tls", "tech", "cookies",
	"links", "scripts", "forms", "robots", "sitemap",
	"title", "favicon",
}

func main() {
	domainPtr := flag.String("d", "", "Target domain (e.g. example.com)")
	funcPtr := flag.String("f", "", "Module or Profile to run (e.g. dns, full, quick)")
	jsonPtr := flag.Bool("json", false, "Output results as JSON")
	outPtr := flag.String("o", "", "Directory to save logs/results")
	versionPtr := flag.Bool("version", false, "Show version")

	flag.Usage = func() {
		output.PrintBanner(version)
		fmt.Fprintf(os.Stderr, "Usage: netscope -d <target> -f <mode>\n\n")

		fmt.Fprintf(os.Stderr, "=== Quick Start ===\n")
		fmt.Fprintf(os.Stderr, "  netscope -d example.com -f quick\n")
		fmt.Fprintf(os.Stderr, "  netscope -d example.com -f full --json -o ./results\n\n")

		fmt.Fprintf(os.Stderr, "=== Profiles ===\n")
		for k, m := range profileMap {
			fmt.Fprintf(os.Stderr, "  %-12s : %s\n", k, strings.Join(m, ", "))
		}
		
		fmt.Fprintf(os.Stderr, "\n=== Modules ===\n")
		fmt.Fprintf(os.Stderr, "  %s\n\n", strings.Join(allModules, ", "))

		fmt.Fprintf(os.Stderr, "=== Special ===\n")
		fmt.Fprintf(os.Stderr, "  netscope doctor  : Run environment checks\n")
	}

	// Check if doctor
	if len(os.Args) > 1 && os.Args[1] == "doctor" {
		output.PrintBanner(version)
		fmt.Println("Running Doctor Checks...")
		res := utils.RunDoctor()
		if res.Internet {
			fmt.Println("✅ Internet : PASS")
		} else {
			fmt.Println("❌ Internet : FAIL")
		}
		if res.DNS {
			fmt.Println("✅ DNS      : PASS")
		} else {
			fmt.Println("❌ DNS      : FAIL")
		}
		for _, iss := range res.Issues {
			fmt.Println("   - Issue: " + iss)
		}
		os.Exit(0)
	}

	flag.Parse()

	if *versionPtr {
		output.PrintBanner(version)
		os.Exit(0)
	}

	if *domainPtr == "" || *funcPtr == "" {
		flag.Usage()
		os.Exit(1)
	}

	domain := utils.NormalizeDomain(*domainPtr)
	if !*jsonPtr {
		output.PrintBanner(version)
	}

	// Determine modules to run
	var toRun []string
	if val, ok := profileMap[*funcPtr]; ok {
		toRun = val
	} else {
		toRun = []string{*funcPtr} // run as single module
	}

	ctx := context.Background()
	results := make(map[string]interface{})

	// Group modules
	var webMods []string
	for _, m := range toRun {
		switch m {
		case "dns":
			results["dns"] = dns.GetDNS(ctx, domain, 5*time.Second)
		case "whois":
			results["whois"] = dns.GetWhois(domain)
		case "ip":
			results["ip"] = dns.GetIP(domain)
		case "subs":
			results["subs"] = dns.GetSubdomains(domain)
		case "email":
			results["email"] = dns.GetEmailSecurity(ctx, domain, 5*time.Second)
		case "status", "headers", "redirect", "tls", "tech", "cookies", "links", "scripts", "forms", "robots", "sitemap", "title", "favicon", "ports":
			webMods = append(webMods, m)
		default:
			fmt.Fprintf(os.Stderr, "[!] Unknown module or profile: %s\n", m)
		}
	}

	if len(webMods) > 0 {
		webRes := web.AnalyzeWeb(domain, webMods)
		for _, wm := range webMods {
			switch wm {
			case "status":
				results["status"] = webRes.Status
			case "headers":
				results["headers"] = webRes.Headers
			case "redirect":
				results["redirect"] = webRes.Redirect
			case "tls":
				results["tls"] = webRes.TLS
			case "cookies":
				results["cookies"] = webRes.Cookies
			case "tech":
				results["tech"] = webRes.Tech
			case "ports":
				results["ports"] = webRes.Ports
			case "robots":
				results["robots"] = webRes.Robots
			case "sitemap":
				results["sitemap"] = webRes.Sitemap
			case "title":
				results["title"] = webRes.Title
			case "favicon":
				results["favicon"] = webRes.Favicon
			case "links":
				if webRes.Crawl != nil {
					results["links"] = webRes.Crawl.Links
				}
			case "scripts":
				if webRes.Crawl != nil {
					results["scripts"] = webRes.Crawl.Scripts
				}
			case "forms":
				if webRes.Crawl != nil {
					results["forms"] = webRes.Crawl.Forms
				}
			}
		}
		if webRes.Error != "" {
			results["web_error"] = webRes.Error
		}
	}

	output.PrintResults(results, *jsonPtr, *outPtr, domain)
}
