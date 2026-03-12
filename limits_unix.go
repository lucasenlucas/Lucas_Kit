//go:build !windows

package main

import (
	"fmt"
	"syscall"
)

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
