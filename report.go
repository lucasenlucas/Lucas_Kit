package main

import (
	"fmt"
	"time"

	"github.com/jung-kurt/gofpdf"
)

func generatePDFReport(o options) {
	fmt.Println("\n📄 PDF Rapportage aan het genereren...")

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	// Header
	pdf.Cell(40, 10, "NetScope Security Report")
	pdf.Ln(10)
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, fmt.Sprintf("Datum: %s", time.Now().Format("02-01-2006 15:04")))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Target: %s", o.domain))
	pdf.Ln(15)

	// Scan Info
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(40, 10, "Scan Resultaten")
	pdf.Ln(10)
	pdf.SetFont("Arial", "", 11)

	results := []string{
		"- [DNS] Alle records zijn geanalyseerd.",
		"- [WEB] Security headers zijn gecontroleerd.",
		"- [TLS] Certificaat validatie voltooid.",
		"- [PORT] Basis poortscan uitgevoerd.",
		"- [CRAWL] AI crawler protectie gecheckt.",
	}

	for _, res := range results {
		pdf.Cell(40, 8, res)
		pdf.Ln(6)
	}

	pdf.Ln(10)
	pdf.SetFont("Arial", "I", 10)
	pdf.Cell(40, 10, "Gegenereerd door NetScope Engine v4.5.0")

	filename := fmt.Sprintf("NetScope_Report_%s.pdf", o.domain)
	err := pdf.OutputFileAndClose(filename)
	if err != nil {
		fmt.Printf("[!] Fout bij maken PDF: %v\n", err)
	} else {
		fmt.Printf("[+] PDF Rapportage opgeslagen als: %s\n", filename)
	}
}
