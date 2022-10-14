package main

import (
	"os"

	"github.com/Digital-Voting-Team/auth-serivce/internal/cli"
)

func main() {
	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}
