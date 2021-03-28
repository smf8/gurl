package main

import (
	"github.com/smf8/gurl/internal/app/gurl/cmd"
	"os"
)

func main() {
	root := cmd.NewCommand()

	if root != nil {
		if err := root.Execute(); err != nil {
			os.Exit(1)
		}
	}
}
