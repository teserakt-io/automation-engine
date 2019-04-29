package main

import (
	"fmt"
	"log"

	"gitlab.com/teserakt/c2se/internal/cli"
	"gitlab.com/teserakt/c2se/internal/cli/commands"
)

// Provided by build script
var gitCommit string
var gitTag string
var buildDate string

func main() {

	log.SetFlags(0)

	clientFactory := cli.NewAPIClientFactory()

	rootCmd := commands.NewRootCommand(clientFactory, getVersion())
	if err := rootCmd.CobraCmd().Execute(); err != nil {
		log.Fatal(err)
	}
}

func getVersion() string {
	var out string

	if len(gitTag) == 0 {
		out = fmt.Sprintf("E4: C2 script reader cli - version %s-%s\n", buildDate, gitCommit)
	} else {
		out = fmt.Sprintf("E4: C2 script reader cli - version %s (%s-%s)\n", gitTag, buildDate, gitCommit)
	}
	out += fmt.Sprintln("Copyright (c) Teserakt AG, 2018-2019")

	return out
}
