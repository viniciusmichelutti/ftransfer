package main

import (
	"fmt"
	"os"

	"github.com/viniciusmichelutti/ftransfer/internal/adapter/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "ftransfer:", err)
		os.Exit(1)
	}
}
