// Command cinegraph starts the CineGraph MCP server.
package main

import (
	"flag"
	"fmt"
	"os"
)

var version = "dev"

func main() {
	showVersion := flag.Bool("version", false, "print the CineGraph version")
	flag.Parse()

	if *showVersion {
		fmt.Println(version)
		return
	}

	fmt.Fprintln(os.Stderr, "CineGraph MCP server is not implemented yet")
	os.Exit(1)
}
