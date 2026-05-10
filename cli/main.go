package main

import (
	"fmt"
	"os"

	"github.com/Pumahawk/simpl-agents-controller/internal/cmd"
	"github.com/Pumahawk/simpl-agents-controller/internal/gitlab"
)

var glclient = gitlab.Client{
	BaseUrl: "https://code.europa.eu",
}

var cmds cmd.Command = &cmd.Group{
	Cmds: []cmd.Command{},
}

func main() {
	args := os.Args
	if err := cmds.Run(args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
