# NetScope (v5.0.0) ⚡️ GOD MODE
**The Ultimate DNS, Web Analysis & Stress Testing Engine.**

NetScope is a powerful, flag-based CLI tool designed for comprehensive domain analysis and network stress testing. Formerly known as Lucas Kit, it has evolved into a streamlined platform for security researchers and system administrators.

**Author**: Lucas Mangroelal | [lucasmangroelal.nl](https://lucasmangroelal.nl)

---

## 🎯 The Concept: Flag-Based Architecture
NetScope v4.7.0 introduces a new, simplified command structure. Instead of dozens of individual flags, we use a unified system:

- **`-d {domein}`**: Het doelwit dat je wilt analyseren (bijv. `example.com`).
- **`-f {functie}`**: De specifieke analyse of actie die je wilt uitvoeren.

Dit maakt de tool sneller, overzichtelijker en makkelijker te combineren met andere scripts.

---

## 🔍 Beschikbare Functies (`-f`)
Je kunt kiezen uit **20 krachtige functies** verdeeld over verschillende categorieën:

### 🌐 Discovery & DNS
- `subs`: Zoekt subdomeinen via publieke certificaatlogs (CT-logs).
- `dns`: Toont alle belangrijke records (A, AAAA, MX, NS, TXT).
- `whois`: Registratiegegevens en nameservers.
- `ip`: Directe lookup van het IP-adres.

### 🛡️ Web & Security Analysis
- `status`: Haalt de HTTP status code op.
- `headers`: Controleert op security headers (HSTS, CSP, X-Frame, etc.).
- `ssl`: Diepgaande check van het SSL/TLS certificaat.
- `redirect`: Volgt en toont de volledige redirect-keten.
- `tech`: Identificeert CMS (zoals WordPress), serversoftware en frameworks.
- `cookies`: Analyseert cookie-attributen (Secure, HttpOnly, SameSite).
- `ports`: Scant veelvoorkomende open webpoorten.
- `email`: Checkt mailbeveiliging (SPF, DKIM, DMARC records).

### 📂 Content & Resources
- `robots`: Bekijkt de inhoud van `robots.txt`.
- `sitemap`: Zoekt naar de `sitemap.xml`.
- `links`: Extraheert alle interne en externe links van de pagina.
- `scripts`: Lijst alle geladen externe JavaScript bestanden op.
- `forms`: Vindt alle HTML-formulieren en hun verzendmethoden.
- `favicon`: Haalt de URL van de favicon op.
- `title`: Geeft de paginatitel (`<title>`) weer.

---

## 🔥 Interactive Attack Wizard (`-f attack`)
De aanval-module is nu volledig interactief. Wanneer je `netscope -d target.com -f attack` runt:
1. **Meting**: De tool meet eerst zelf de snelheid van de site.
2. **Advies**: Je krijgt een geadviseerd "Attack Level" gebaseerd op de latency.
3. **Configuratie**: Kies je gewenste level (1-10) en de tijdsduur.
4. **Monitoring**: De tool houdt live bij of de site nog online is of plat gaat.

---

## 🛠️ Installatie

### Snel Installeren (macOS & Linux)
```bash
curl -fsSL https://raw.githubusercontent.com/lucasenlucas/NetScope/main/scripts/install.sh | sh
```

### Snel Installeren (Windows)
```powershell
irm https://raw.githubusercontent.com/lucasenlucas/NetScope/main/scripts/install.ps1 | iex
```

### Handmatige Build (Source)
```bash
git clone https://github.com/lucasenlucas/NetScope.git
cd NetScope
go build -o netscope
./netscope --help
```

---

## 💡 Voorbeelden
**DNS Analyse uitvoeren:**
```bash
netscope -d google.com -f dns
```

**Security Headers checken:**
```bash
netscope -d nu.nl -f headers
```

**Een stresstest starten:**
```bash
netscope -d eigen-test-site.nl -f attack
```

---
_NetScope is onderdeel van gost three tooling by Lucas Mangroelal and Quin de Lira._
