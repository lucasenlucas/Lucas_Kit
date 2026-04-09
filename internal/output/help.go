package output

import (
	"fmt"
	"os"
	"strings"
)

var descriptions = map[string]struct {
	short string
	long  string
	usage string
}{
	// Profiles
	"quick": {
		short: "Fast first impression of a domain.",
		long:  "The 'quick' profile is designed to give you an immediate overview of the target without generating heavy traffic. It resolves the IP, checks if the server is up (HTTP Status), fetches basic Security Headers, verifies the TLS/SSL certificate, and gets the Page Title.",
		usage: "netscope -d example.com -f quick",
	},
	"web": {
		short: "Deep web component analysis.",
		long:  "The 'web' profile simulates a deep scanner checking for redirects, security headers, tech stacks, available cookies, and statically extracting scripts/links/forms from the HTML body.",
		usage: "netscope -d example.com -f web",
	},
	"dns-full": {
		short: "Infrastructure & DNS intelligence.",
		long:  "The 'dns-full' profile thoroughly maps out the underlying infrastructure. It grabs WHOIS information, all relevant DNS records (A, AAAA, MX, NS, TXT), email security configurations (SPF, DKIM, DMARC), and searches for known subdomains via crt.sh.",
		usage: "netscope -d example.com -f dns-full",
	},
	"full": {
		short: "Complete domain audit.",
		long:  "The 'full' profile runs every single module NetScope has to offer. Be prepared for a longer execution time. It covers web, infrastructure, crawls, and network scanning all at once.",
		usage: "netscope -d example.com -f full --json -o ./results",
	},
	"pageinfo": {
		short: "Retrieves title and favicon.",
		long:  "Scrapes the primary index page for its <title> tag and attempts to locate the favicon URL.",
		usage: "netscope -d example.com -f pageinfo",
	},
	"web-basic": {
		short: "Basic HTTP headers and status.",
		long:  "Returns the HTTP status code, security headers, and the redirect chain if the server forwards you to another URL.",
		usage: "netscope -d example.com -f web-basic",
	},
	"crawl-lite": {
		short: "Statically extracts links, scripts, and forms.",
		long:  "Downloads the HTML body and statically extracts all <a> hrefs, <script> sources, and <form> action URIs.",
		usage: "netscope -d example.com -f crawl-lite",
	},

	// Modules
	"dns": {
		short: "Queries A, AAAA, MX, NS, TXT records.",
		long:  "Performs standard DNS queries against your local resolver to grab standard internet routing records.",
		usage: "netscope -d example.com -f dns",
	},
	"whois": {
		short: "Grabs authoritative nameservers.",
		long:  "A lightweight lookup to find the authoritative nameservers for the domain.",
		usage: "netscope -d example.com -f whois",
	},
	"ip": {
		short: "Resolves domain to IP addresses.",
		long:  "Finds the IPv4/IPv6 addresses bound to the domain.",
		usage: "netscope -d example.com -f ip",
	},
	"subs": {
		short: "Finds subdomains via certificate transparency.",
		long:  "Queries crt.sh to find subdomains that have had TLS certificates issued.",
		usage: "netscope -d example.com -f subs",
	},
	"email": {
		short: "Checks SPF, DKIM, DMARC, MTA-STS.",
		long:  "Checks the domain's TXT records for email spoofing protection protocols like SPF and DMARC.",
		usage: "netscope -d example.com -f email",
	},
	"ports": {
		short: "Checks common web ports (80, 443, 8080, 8443).",
		long:  "Attempts a quick 1-second TCP handshake on common web ports to see if they are open.",
		usage: "netscope -d example.com -f ports",
	},
	"status": {
		short: "Returns HTTP status code.",
		long:  "Sends an HTTP GET request and returns the resulting status code (e.g. 200, 404, 503).",
		usage: "netscope -d example.com -f status",
	},
	"headers": {
		short: "Checks for security headers.",
		long:  "Inspects the HTTP response for best-practice security headers like Content-Security-Policy or Strict-Transport-Security.",
		usage: "netscope -d example.com -f headers",
	},
	"redirect": {
		short: "Traces the HTTP redirect chain.",
		long:  "Follows up to 10 HTTP redirects and maps out the exact URLs the server routes you through.",
		usage: "netscope -d example.com -f redirect",
	},
	"tls": {
		short: "Extracts TLS/SSL certificate details.",
		long:  "Establishes a secure connection and extracts the Certificate Issuer, Expiry Date, and valid DNS Subject Alternative Names.",
		usage: "netscope -d example.com -f tls",
	},
	"tech": {
		short: "Fingerprints server technologies.",
		long:  "Looks at HTTP headers and HTML content to identify backend software like WordPress, React, Apache, Nginx, etc.",
		usage: "netscope -d example.com -f tech",
	},
	"cookies": {
		short: "Analyzes set cookies for security flags.",
		long:  "Checks cookies given by the server for HttpOnly, Secure, and SameSite attributes.",
		usage: "netscope -d example.com -f cookies",
	},
	"links": {
		short: "Extracts all <a> tags from the HTML.",
		long:  "Parses the body of the response to find any hyperlinked URLs on the start page.",
		usage: "netscope -d example.com -f links",
	},
	"scripts": {
		short: "Extracts all <script> sources.",
		long:  "Lists all external or bundled JS files injected into the page's HTML.",
		usage: "netscope -d example.com -f scripts",
	},
	"forms": {
		short: "Extracts <form> actions and methods.",
		long:  "Useful for finding login pages, search inputs, or unprotected submission endpoints on the page.",
		usage: "netscope -d example.com -f forms",
	},
	"robots": {
		short: "Fetches robots.txt.",
		long:  "Downloads and displays the contents of the robots.txt file if it exists, showing disallowed crawler paths.",
		usage: "netscope -d example.com -f robots",
	},
	"sitemap": {
		short: "Checks for sitemap.xml.",
		long:  "Does a quick check to see if /sitemap.xml is accessible on the server.",
		usage: "netscope -d example.com -f sitemap",
	},
	"title": {
		short: "Extracts the page <title>.",
		long:  "Grabs the HTML title element, showing exactly what the page claims to be in browser tabs.",
		usage: "netscope -d example.com -f title",
	},
	"favicon": {
		short: "Locates the favicon URL.",
		long:  "Scans for a favicon link tag in the HTML or checks the standard /favicon.ico path.",
		usage: "netscope -d example.com -f favicon",
	},
}

var allProfiles = []string{"quick", "web", "dns-full", "full", "pageinfo", "web-basic", "crawl-lite"}
var allModules = []string{
	"dns", "whois", "ip", "subs", "email", "ports",
	"status", "headers", "redirect", "tls", "tech", "cookies",
	"links", "scripts", "forms", "robots", "sitemap",
	"title", "favicon",
}

// PrintGeneralHelp prints a clear and organized overall help menu
func PrintGeneralHelp(version string) {
	PrintBanner(version)
	fmt.Fprintf(os.Stderr, "NetScope is a fast domain intelligence CLI. Gather recon, web structure, and security insight in a single command.\n\n")

	fmt.Fprintf(os.Stderr, "🚀 USAGE\n")
	fmt.Fprintf(os.Stderr, "  netscope -d <domain> -f <mode> [options]\n\n")

	fmt.Fprintf(os.Stderr, "🎯 PROFILES (Groups of modules)\n")
	for _, p := range allProfiles {
		fmt.Fprintf(os.Stderr, "  %-14s %s\n", p, descriptions[p].short)
	}

	fmt.Fprintf(os.Stderr, "\n🧩 MODULES (Run individually)\n")
	for i, m := range allModules {
		fmt.Fprintf(os.Stderr, "  %-14s", m)
		if (i+1)%4 == 0 {
			fmt.Fprintf(os.Stderr, "\n")
		}
	}
	fmt.Fprintf(os.Stderr, "\n\n⚙️  FLAGS\n")
	fmt.Fprintf(os.Stderr, "  -d, --domain   Target domain to analyze (e.g., example.com)\n")
	fmt.Fprintf(os.Stderr, "  -f, --flag     Profile or module to execute (e.g., quick, dns)\n")
	fmt.Fprintf(os.Stderr, "  --json         Format the output strictly as JSON\n")
	fmt.Fprintf(os.Stderr, "  -o             Directory to save the JSON output file\n")
	fmt.Fprintf(os.Stderr, "  doctor         Run NetScope system environment checking\n")
	
	fmt.Fprintf(os.Stderr, "\n💡 TIP: For detailed info on any profile or module, run:\n")
	fmt.Fprintf(os.Stderr, "  netscope --help <name>  (For example: netscope --help tls)\n\n")
}

// PrintDetailedHelp prints deep information about a specific command/flag
func PrintDetailedHelp(version, command string) {
	cmdLower := strings.ToLower(command)
	info, exists := descriptions[cmdLower]
	
	if !exists {
		PrintBanner(version)
		fmt.Fprintf(os.Stderr, "❌ The module or profile '%s' does not exist.\n", command)
		fmt.Fprintf(os.Stderr, "Run 'netscope --help' to see all available options.\n\n")
		return
	}

	PrintBanner(version)
	fmt.Fprintf(os.Stderr, "📖 DETAILED HELP: %s\n", strings.ToUpper(cmdLower))
	fmt.Fprintf(os.Stderr, strings.Repeat("=", 40)+"\n\n")
	
	fmt.Fprintf(os.Stderr, "Description:\n")
	
	// Create line breaks nicely for terminal
	words := strings.Split(info.long, " ")
	line := "  "
	for _, w := range words {
		if len(line)+len(w) > 80 {
			fmt.Fprintf(os.Stderr, "%s\n", line)
			line = "  " + w + " "
		} else {
			line += w + " "
		}
	}
	if strings.TrimSpace(line) != "" {
		fmt.Fprintf(os.Stderr, "%s\n", line)
	}

	fmt.Fprintf(os.Stderr, "\nExample Usage:\n")
	fmt.Fprintf(os.Stderr, "  $ %s\n", info.usage)
	
	fmt.Fprintf(os.Stderr, "\nOther Options:\n")
	fmt.Fprintf(os.Stderr, "  Combine with '--json' for pipeline parsing.\n")
	fmt.Fprintf(os.Stderr, "  Combine with '-o ./outputs' to securely save the generated payload.\n\n")
}
