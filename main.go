package main

import (
	"github.com/cloudfauj/cloudfauj/cmd"
)

func main() {
	// TODO: Create & supply api client from here instead of
	//  creating them in the specific sub-commands.
	cmd.Execute()
}
