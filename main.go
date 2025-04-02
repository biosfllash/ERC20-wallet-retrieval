package main

import (
	"fmt"
	"walletDev/gui"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	// Build information for the application.
)

func main() {
	fmt.Printf("Wallet Manager %s\nCommit: %s\nBuilt: %s\n",
		version, commit, date)
	gui.Start()
}
