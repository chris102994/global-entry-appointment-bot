package main

import "github.com/chris102994/global-entry-appointment-bot/cmd"

var (
	Branch         = "N/A"
	BuildTimestamp = "N/A"
	CommitHash     = "N/A"
	Version        = "N/A"
)

func main() {
	if err := cmd.NewRootCmd(Branch, BuildTimestamp, CommitHash, Version).Execute(); err != nil {
		panic(err)
	}
}
