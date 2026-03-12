package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

const version = "4.8.0"

type options struct {
	domain   string
	function string
	help     bool
	version  bool

	// General Options
	jsonOut   bool
	outputDir string
	check     bool
	update    bool

	// Stress Test internal config (set via attack level)
	probes        int
	attackMinutes int
	concurrency   int
	level         int
	noKeepAlive   bool

	// DNS config
	resolver string
	timeout  time.Duration

	// Proxy support
	proxyFile string
}

func main() {
	var o options

	// General
	flag.StringVar(&o.domain, "d", "", "Target domain (e.g. example.com)")
	flag.StringVar(&o.function, "f", "", "Function to run (e.g. dns, ssl, headers)")
	flag.BoolVar(&o.jsonOut, "json", false, "Output results as JSON")
	flag.StringVar(&o.outputDir, "o", "", "Directory to save logs/results")
	flag.BoolVar(&o.help, "help", false, "Show help page")
	flag.BoolVar(&o.help, "h", false, "Short help flag")
	flag.BoolVar(&o.version, "version", false, "Show version")
	flag.BoolVar(&o.check, "check", false, "Check for updates")
	flag.BoolVar(&o.update, "update", false, "Update to latest version")
	flag.StringVar(&o.proxyFile, "proxies", "", "File with proxy list (IP:Port)")

	// Internal Stress/DNS defaults
	o.timeout = 5 * time.Second
	o.probes = 1

	flag.Usage = func() {
		printBanner(version)
		fmt.Fprintf(os.Stderr, "Usage: netscope -d <target> -f <function>\n\n")

		printBoxedSection("🎯 CORE COMMANDS", []flagHelp{
			{"-d, --domain", "Het domein dat je wilt analyseren of testen"},
			{"-f, --flag", "De specifieke functie die je wilt uitvoeren"},
			{"--help", "Toon dit help-scherm"},
			{"--version", "Toon huidige versie"},
		})

		printBoxedSection("🔍 ANALYSIS FUNCTIONS", []flagHelp{
			{"subs", "Zoekt bekende subdomeinen via certificaatlogs"},
			{"dns", "Toont A, AAAA, MX, NS en TXT records"},
			{"whois", "Toont registrar en registratie info"},
			{"ssl", "Controleert het SSL-certificaat"},
			{"headers", "Haalt HTTP security headers op"},
			{"status", "Geeft de HTTP status code terug"},
			{"redirect", "Toont de volledige redirect-keten"},
			{"robots", "Haalt robots.txt bestand op"},
			{"sitemap", "Zoekt naar sitemap.xml"},
			{"tech", "Identificeert CMS, frameworks en server software"},
			{"forms", "Zoekt formulieren op de pagina"},
			{"links", "Haalt alle links van de pagina op"},
			{"scripts", "Toont alle geladen externe scripts"},
			{"cookies", "Analyseert cookies (Secure, HttpOnly, SameSite)"},
			{"email", "Checkt mailbeveiliging (SPF, DKIM, DMARC)"},
			{"ports", "Checkt veelvoorkomende open webpoorten"},
			{"favicon", "Haalt de favicon van de site op"},
			{"title", "Toont de paginatitel"},
			{"ip", "Zoekt het IP-adres van het domein op"},
			{"attack", "Interactieve stresstest wizard"},
		})

		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  netscope -d example.com -f dns\n")
		fmt.Fprintf(os.Stderr, "  netscope -d example.com -f ssl\n")
		fmt.Fprintf(os.Stderr, "  netscope -d example.com -f attack\n")
	}

	flag.Parse()

	if o.version {
		printBanner(version)
		os.Exit(0)
	}

	if o.help || (len(os.Args) == 1) {
		flag.Usage()
		os.Exit(0)
	}

	if o.update {
		runAutoUpdate()
		os.Exit(0)
	}

	if o.check {
		runCheckUpdate()
		os.Exit(0)
	}

	if o.domain == "" {
		fmt.Println("[!] Domein (-d) is verplicht.")
		os.Exit(1)
	}

	if o.function == "" {
		fmt.Println("[!] Functie (-f) is verplicht. Gebruik --help voor een lijst met functies.")
		os.Exit(1)
	}

	if !o.jsonOut {
		printBanner(version)
	}

	runFocusedAnalysis(o)
}

func runFocusedAnalysis(o options) {
	ctx := context.Background()
	domain := normalizeDomain(o.domain)

	switch strings.ToLower(o.function) {
	case "subs":
		runSubdomainScan(domain, o)
	case "dns":
		runDNSAnalysis(ctx, o)
	case "whois":
		runWhoisAnalysis(domain, o)
	case "ssl":
		o.function = "ssl" // ensure it's set for web_analysis
		runWebAnalysis(o)
	case "headers":
		runWebAnalysis(o)
	case "status":
		runWebAnalysis(o)
	case "redirect":
		runWebAnalysis(o)
	case "robots":
		runWebAnalysis(o)
	case "sitemap":
		runWebAnalysis(o)
	case "tech":
		runWebAnalysis(o)
	case "forms":
		runWebAnalysis(o)
	case "links":
		runWebAnalysis(o)
	case "scripts":
		runWebAnalysis(o)
	case "cookies":
		runWebAnalysis(o)
	case "email":
		runMailSecurity(ctx, domain, o)
	case "ports":
		runWebAnalysis(o)
	case "favicon":
		runWebAnalysis(o)
	case "title":
		runWebAnalysis(o)
	case "ip":
		runIPLookup(domain, o)
	case "attack":
		runSitePlof(o)
	default:
		fmt.Printf("[!] Onbekende functie: %s\n", o.function)
		fmt.Println("Gebruik --help voor een lijst met alle functies.")
	}
}
