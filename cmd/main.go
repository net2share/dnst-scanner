package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "scan":
		runScan(os.Args[2:]) // <-- EINZIGER Dispatcher
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", os.Args[1])
		usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Println(`Usage:
  dnst-scanner scan --tunnel-domain <domain> [flags]

Commands:
  scan    Run basic resolver scan`)
}
