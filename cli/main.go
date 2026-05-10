package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/Pumahawk/simpl-agents-controller/internal/cmd"
	"github.com/Pumahawk/simpl-agents-controller/internal/gitlab"
)

var glclient = gitlab.Client{
	BaseUrl: "https://code.europa.eu",
}

var cmds cmd.Command = &cmd.Group{
	Cmds: []cmd.Command{
		&RegistryCmd,
	},
}

func main() {
	args := os.Args[1:]
	if err := cmds.Run(args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

var RegistryCmd = cmd.Cmd{
	CName: "registry",
	CRun: func(args []string) error {
		var page, perpage int
		fs := flag.NewFlagSet("", flag.ExitOnError)
		fs.IntVar(&page, "page", 1, "")
		fs.IntVar(&perpage, "perpage", 10, "")
		fs.Parse(args)

		ids := fs.Arg(0)
		if ids == "" {
			return fmt.Errorf("missing project id")
		}

		id, err := strconv.ParseInt(ids, 10, 64)
		if err != nil {
			return err
		}

		items, err := glclient.Packages(int(id), page, perpage)
		if err != nil {
			return err
		}

		w := tabwriter.NewWriter(os.Stdout, 1, 2, 1, ' ', 0)
		fmt.Fprintf(w, "Id\tName\tVersion\tType\tPipelineId\tRef\tPipeline WebUrl\n")
		for i := range items {
			item := items[i]
			fmt.Fprintf(w, "%d \t %q \t %q \t %q \t %d \t %q \t %q\n",
				item.Id, item.Name, item.Version, item.PackageType, item.Pipeline.Id, item.Pipeline.Ref, item.Pipeline.WebUrl)
		}
		w.Flush()
		return nil
	},
}
