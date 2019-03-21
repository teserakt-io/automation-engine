package main

import "fmt"

var gitCommit string
var gitTag string
var buildDate string

func main() {

	if len(gitTag) == 0 {
		fmt.Printf("E4: C2 scripting engine - version %s-%s\n", buildDate, gitCommit[:4])
	} else {
		fmt.Printf("E4: C2 scripting engine - version %s (%s-%s)\n", gitTag, buildDate, gitCommit[:4])
	}
	fmt.Println("Copyright (c) Teserakt AG, 2018-2019")

	// TODO: everything.
}
