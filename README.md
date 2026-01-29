# Lucas Kit (`lucaskit`)

> The ultimate domain toolkit containing **LucasDNS** and **Lucaskill**. Made by Lucas Mangroelal | lucasmangroelal.nl

**Lucas Kit** is een collectie van krachtige tools voor DNS/Domain information gathering en stress testing. Het bevat:

1. **LucasDNS**: Info gathering (DNS, WHOIS, Mail Security, Subdomains).
2. **Lucaskill**: Advanced HTTP stress test / load test tool.

## Install

### Kali Linux / macOS / Linux (aanbevolen)

**Automatische installatie (detecteert architecture):**

```bash
curl -fsSL https://raw.githubusercontent.com/lucasenlucas/Lucas_Kit/main/scripts/install.sh | sh
```

Dit installeert zowel `lucasdns` als `lucaskill` naar `/usr/local/bin` (of `~/.local/bin`).

### Windows

**PowerShell:**

```powershell
.\scripts\install.ps1 -Repo "lucasenlucas/Lucas_Kit"
```

## Tools

### 1. LucasDNS (`lucasdns`)

Info gathering tool.

```bash
lucasdns -d <domein> [flags]
```

**Features:**
- DNS Records (A, AAAA, MX, NS, TXT, SOA, CAA, SRV)
- Mail Security (SPF, DMARC, DKIM, MTA-STS)
- WHOIS informatie
- Certificate Transparency (Subdomeinen)

**Voorbeelden:**
```bash
lucasdns -d example.com -inf -n
lucasdns -d example.com -subs
```

### 2. Lucaskill (`lucaskill`)

HTTP stress/load test tool.

```bash
lucaskill -d <domein> -t <minuten> [flags]
```

**Features:**
- High performance HTTP flooding
- "Keep-Down" logic: monitort site en valt automatisch weer aan als hij online komt
- Real-time dashboard
- Reporting (`-o output_dir`)

**Voorbeelden:**
```bash
lucaskill -d example.com -t 10
```

> **⚠️ DISCLAIMER:** Gebruik Lucaskill alleen op systemen waar je expliciete toestemming voor hebt.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

Additional Notice on Naming, Forks, and Liability

Use of the names "Lucas Mangroelal", "Lucas DNS", "Lucas Kit", or any related project names associated with the original version of this Software does not imply endorsement by the original author.
Any redistributed, modified, or forked versions must make it clear that they are unofficial versions if they are not directly maintained by Lucas Mangroelal.
Lucas Mangroelal is not responsible or liable for any misuse, damages, or consequences resulting from third-party copies, forks, or modified versions of this Software.

For more information, permissions regarding naming, or official inquiries, contact:
kit@lucasmangroelal.nl
