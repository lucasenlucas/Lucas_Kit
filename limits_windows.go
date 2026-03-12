//go:build windows

package main

func setRlimits() {
	// No-op for Windows as Getrlimit/Setrlimit are not available
	// Windows manages handles differently
}
