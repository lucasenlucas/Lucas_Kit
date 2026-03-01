package main

import (
	"context"
	"fmt"
)

func runUnifiedAnalysis(o options) {
	if !o.jsonOut {
		fmt.Printf("\n🚀 NetScope Unified Analysis: %s\n", o.domain)
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	}

	ctx := context.Background()

	// DNS & Mail Analyser routing
	runDNS := false
	if o.inf || o.n || o.whois || o.subs || o.a || o.aaaa || o.cname || o.mx || o.ns || o.txt || o.soa || o.caa || o.srv || o.records != "" || o.resolve != "" || o.dnssec {
		runDNS = true
	}

	if runDNS {
		if !o.jsonOut {
			fmt.Println("\n📡 [MODULE: DNS & MAIL SECURITY]")
		}
		runDNSAnalysis(ctx, o)
	}

	// Web Security routing
	runWeb := false
	if o.httpCheck || o.tlsCheck || o.headersCheck || o.cacheCheck || o.fingerCheck || o.portsCheck || o.pathsCheck || o.corsCheck || o.cookieCheck || o.techCheck || o.crawlerCheck || o.methodCheck {
		runWeb = true
	}

	if runWeb {
		if !o.jsonOut {
			fmt.Println("\n🛡️  [MODULE: WEB SECURITY & ANALYSIS]")
		}
		runWebAnalysis(o)
	}

	// Vulnerability & Discovery routing
	runDiscovery := false
	if o.dirCheck || o.paramsCheck || o.cmsCheck {
		runDiscovery = true
	}
	if runDiscovery {
		if !o.jsonOut {
			fmt.Println("\n🔍 [MODULE: VULNERABILITY & DISCOVERY]")
		}
		runVulnAnalysis(o)
	}

	// Metrics / Measure routing
	if o.measure {
		if !o.jsonOut {
			fmt.Println("\n⚡ [MODULE: L7 METRICS & MEASURE]")
		}
		if !runWeb {
			runWebAnalysis(o)
		}
	}

	if !o.jsonOut {
		fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Println("✅ Alle aangevraagde analyses zijn voltooid.")
	}
}
