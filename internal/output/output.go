package output

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// PrintBanner shows the application banner.
func PrintBanner(version string) {
	banner := `
    _   __     __  _____                     
   / | / /__  / /_/ ___/_________  ____  ___ 
  /  |/ / _ \/ __/\__ \/ ___/ __ \/ __ \/ _ \
 / /|  /  __/ /_ ___/ / /__/ /_/ / /_/ /  __/
/_/ |_/\___/\__//____/\___/\____/ .___/\___/ 
                               /_/           
`
	fmt.Println(banner)
	fmt.Println("NetScope by Lucas Mangroelal | lucasmangroelal.nl")
	fmt.Printf("Version: %s | Platform: %s/%s\n", version, runtime.GOOS, runtime.GOARCH)
	fmt.Println("")
}

// PrintResults renders the structured data.
func PrintResults(data map[string]interface{}, jsonOut bool, outputDir string, domain string) {
	if jsonOut {
		b, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			fmt.Printf("{\"error\": \"Failed to marshal JSON: %v\"}\n", err)
			return
		}
		fmt.Println(string(b))
		
		if outputDir != "" {
			saveToFile(outputDir, fmt.Sprintf("netscope_%s.json", domain), b)
		}
		return
	}

	// Textual render
	fmt.Printf("\n=== Scan Results for %s ===\n", domain)
	for module, result := range data {
		fmt.Printf("\n[ %s ]\n", module)
		printMap(result, "  ")
	}

	if outputDir != "" {
		fmt.Printf("\n[*] Output saving not supported for text mode yet. Use --json to save to file.\n")
	}
}

func printMap(v interface{}, indent string) {
	switch val := v.(type) {
	case map[string]interface{}:
		for k, v2 := range val {
			if isComplex(v2) {
				fmt.Printf("%s- %s:\n", indent, k)
				printMap(v2, indent+"  ")
			} else {
				fmt.Printf("%s- %s: %v\n", indent, k, v2)
			}
		}
	case []string:
		for _, s := range val {
			fmt.Printf("%s* %s\n", indent, s)
		}
	case []interface{}:
		for _, it := range val {
			fmt.Printf("%s* %v\n", indent, it)
		}
	default:
		fmt.Printf("%s%v\n", indent, val)
	}
}

func isComplex(v interface{}) bool {
	switch v.(type) {
	case map[string]interface{}, []string, []interface{}:
		return true
	}
	return false
}

func saveToFile(dir, filename string, data []byte) {
	_ = os.MkdirAll(dir, 0755)
	path := filepath.Join(dir, filename)
	err := os.WriteFile(path, data, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[!] Failed to write to %s: %v\n", path, err)
	} else {
		// Log intentionally minimal. Will be caught by stderr if issues.
	}
}
