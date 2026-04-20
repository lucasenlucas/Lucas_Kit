<p align="center">
  <img src="https://github.com/lucasenlucas/lucas_cdn/blob/main/Scherm%C2%ADafbeelding%202026-04-08%20om%2022.00.14.png?raw=true" alt="NetScope Banner" width="800px"/>
</p>

<p align="center">
  <h1 align="center">NetScope</h1>
  <p align="center">
    <strong>Ultra-Fast Domain Intelligence & Reconnaissance Engine</strong>
    <br />
    Gather deep DNS, web structure, and security insights in a single strike.
  </p>
</p>

<p align="center">
  <a href="https://github.com/lucasenlucas/NetScope/releases"><img src="https://img.shields.io/github/v/release/lucasenlucas/NetScope?style=for-the-badge&color=00FFCC" alt="Release"></a>
  <a href="LICENSE"><img src="https://img.shields.io/github/license/lucasenlucas/NetScope?style=for-the-badge&color=00FFCC" alt="License"></a>
  <a href="https://netseries.dev"><img src="https://img.shields.io/badge/Website-Netseries.dev-00FFCC?style=for-the-badge" alt="Website"></a>
</p>

---

## 🚀 Overview

**NetScope** is a high-performance CLI tool designed for developers, security researchers, and sysadmins who need immediate, structured intelligence on any domain. Forget running ten different tools; NetScope consolidates deep reconnaissance into a single, modular workflow.

### Key Capabilities
- 🛡️ **Security Audit**: Instant detection of security headers, TLS misconfigurations, and cookie flags.
- 🌐 **Web Intelligence**: Analyze tech stacks, redirect chains, and page structure.
- 📡 **Infrastructure Mapping**: Complete DNS record profiling, WHOIS lookups, and subdomain discovery.
- 📧 **Email Security**: Rapid SPF, DKIM, and DMARC verification.

---

## ⚡ Quick Start

### Installation (Mac & Linux)
Install instantly via our optimized installer:
```bash
curl -sL https://raw.githubusercontent.com/lucasenlucas/NetScope/main/install.sh | bash
```

### Installation (Windows)
Run in PowerShell as Administrator:
```powershell
iwr -useb https://raw.githubusercontent.com/lucasenlucas/NetScope/main/scripts/install.ps1 | iex
```

> **Manual Install**: `go install github.com/lucasenlucas/NetScope/cmd/netscope@latest`

---

## 🎯 Usage & Profiles

NetScope uses **Profiles** to group modules for specific use cases.

### 1. Quick Scan
Ideal for a fast first impression.
```bash
netscope -d example.com -f quick
```
*Resolves IP, checks HTTP status, fetches security headers, verifies TLS, and gets page title.*

### 2. Physical Web Analysis
Deep dive into web component behavior.
```bash
netscope -d example.com -f web
```
*Analyzes redirects, cookies, tech stack, and extracts scripts/links/forms.*

### 3. DNS & Infrastructure
Full mapping of the domain's backbone.
```bash
netscope -d example.com -f dns-full
```
*Queries all DNS records, WHOIS, Email security (SPF/DMARC), and subdomains.*

### 4. Full Audit
The complete engine. Everything, all at once.
```bash
netscope -d example.com -f full --json -o ./reports
```

---

## 🧩 Modular Power

Run individual modules for surgical precision:
| Module | Command | Description |
| :--- | :--- | :--- |
| **DNS** | `netscope -d <dom> -f dns` | A, AAAA, MX, NS, TXT |
| **TLS** | `netscope -d <dom> -f tls` | Certificate detail extraction |
| **Tech** | `netscope -d <dom> -f tech` | Fingerprint server stack |
| **Subs** | `netscope -d <dom> -f subs` | CT-log subdomain discovery |
| **Mail** | `netscope -d <dom> -f email` | SPF, DKIM, DMARC check |

---

## 🛠️ Diagnostics & Output

### Doctor Mode
Verify your environment and connectivity:
```bash
netscope doctor
```

### JSON Export
Format output for integration into pipelines:
```bash
netscope -d example.com -f full --json
```

---

## 📦 Project Structure
```text
/cmd/netscope      # Entry point
/internal/         # Core logic
  ├── dns/         # Infrastructure modules
  ├── web/         # HTTP/Scraping modules
  ├── crawl/       # Static extraction
  ├── output/      # Banner & Rendering
  └── utils/       # Shared helpers
```

---

## 💎 Why NetScope?
- **Speed**: Built in Go with concurrency at its core.
- **Clarity**: Beautiful, human-readable terminal output.
- **Precision**: No bloat, just the data you need.
- **Modular**: Easily extensible for custom discovery packs.

---

## 🤝 Support & Contribution
Developed and maintained by the **Netseries Team**.

- 🌐 **Website**: [Netseries.dev](https://netseries.dev)
- 📧 **Contact**: [Team@netseries.dev](mailto:Team@netseries.dev)
- ⭐️ **Star this repo** if you find it useful!
- 🛠️ **Contributions** are welcome via PRs.

---

### ⚠️ Disclaimer
This tool is intended for **authorized testing and educational purposes only**. Always ensure you have permission before scanning any target infrastructure.

---
<p align="center">
  Built with ❤️ by the <strong>Netseries Team</strong>
</p>
