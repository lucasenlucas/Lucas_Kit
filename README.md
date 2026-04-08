<p align="center">
  <img src="https://github.com/lucasenlucas/lucas_cdn/blob/main/Scherm%C2%ADafbeelding%202026-04-08%20om%2022.00.14.png?raw=true" alt="NetScope Banner"/>
</p>

<p align="center">
  Fast Domain Intelligence CLI for Recon, Web Analysis & Security Checks
</p>

<p align="center">
  <strong>1 command → full insight into any domain</strong>
</p>

---

## What is NetScope?

NetScope is a powerful CLI tool that allows you to quickly analyze a domain and understand:

-  DNS & infrastructure  
-  Web behavior & security  
-  Technology stack  
-  Page structure & content  
-  Overall domain intelligence  

Built for developers, security enthusiasts and builders who want **fast, clean and structured results**.

---

## ⚡ Quick Start

### Easy Installation (Mac & Linux)
You can install NetScope instantly using our quick-install script:
```bash
curl -sL https://raw.githubusercontent.com/lucasenlucas/NetScope/main/install.sh | bash
```
> Or install from source via Go: `go install github.com/lucasenlucas/NetScope/cmd/netscope@latest`

### Run a Scan

```
netscope -d example.com -f quick
```
```
netscope -d example.com -f web
```
```
netscope -d example.com -f full --json -o ./reports
```

## Scan Profiles
NetScope is built around profiles so you don’t need to remember 20 commands.

### quick
Fast first impression of a domain

Includes:
- IP lookup
- HTTP status
- Security headers
- TLS check
- Page title

### web
Deep web analysis

Includes:
- Headers
- Redirects
- Cookies
- Tech stack detection
- Forms, scripts & links

### dns-full
Infrastructure & DNS intelligence

Includes:
- DNS records
- WHOIS data
- IP resolution
- Mail security (SPF, DKIM, DMARC)
- Subdomains

### full
Complete domain audit
Runs everything and combines all results.

## Modules (Advanced Usage)
You can still run individual modules:
```
netscope -d example.com -f dns
netscope -d example.com -f headers
netscope -d example.com -f tls
netscope -d example.com -f links
```

## Extra Commands

### pageinfo
```
netscope -d example.com -f pageinfo
```
Returns page title + favicon

### web-basic
```
netscope -d example.com -f web-basic
```
Returns status + headers + redirects

### crawl-lite
```
netscope -d example.com -f crawl-lite
```
Extracts:
- Links
- Scripts
- Forms

## Output Options
Export your results:
```
netscope -d example.com -f full --json
```
```
netscope -d example.com -f full -o ./reports
```

## Doctor Mode
Check your environment:
```
netscope doctor
```
Validates:
* Network connectivity
* DNS resolution
* Configuration

## CLI Usage
```
netscope -d <domain> -f <mode>
```
Flags:
* `-d, --domain` → Target domain
* `-f, --flag` → Scan mode
* `--json` → JSON output
* `-o` → Save results
* `--proxies` → Use proxy file
* `--update` → Check for updates

## Project Structure
```
/cmd/netscope
/internal/
  dns/
  web/
  crawl/
  output/
  utils/
```

## Why NetScope?
Because most tools are either:
* too slow
* too complex
* or too fragmented

NetScope gives you:
*  Speed
*  Clarity
*  Modular design
*  Clean output

## Roadmap
*  Improved reporting system
*  Better tech detection
*  Plugin system
*  API mode

## Author
Built by (me) Lucas Mangroelal  https://lucasmangroelal.nl :) 

## Support ❤️
If you like this project:
- Star the repo
- Contribute
- Share it

### Disclaimer ⚠️
```
This tool is intended for educational and authorized testing purposes only. 
Do not use NetScope on systems you do not own or have permission to test.
```
