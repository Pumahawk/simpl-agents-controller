package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
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
		&LastVersionCmd,
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

var LastVersionCmd = cmd.Cmd{
	CName: "last-version",
	CRun: func(args []string) error {
		var ref, typ string
		fs := flag.NewFlagSet("", flag.ExitOnError)
		fs.StringVar(&ref, "ref", "main", "")
		fs.StringVar(&typ, "type", "helm", "")
		fs.Parse(args)

		projectIds := fs.Args()

		if len(projectIds) == 0 {
			return fmt.Errorf("missing project id")
		}

		nums, err := toInt(projectIds)
		if err != nil {
			return err
		}

		cv := GetVersions(nums)
		w := tabwriter.NewWriter(os.Stdout, 1, 2, 1, ' ', 0)
		find := make(map[int]bool)
		fmt.Fprintf(w, "Project\tVersion\n")
		for v := range cv {
			if v.Ref == ref && v.Type == typ && !strings.Contains(v.Version, "latest") {
				if !find[v.ProjectId] {
					fmt.Fprintf(w, "%d\t%q\n", v.ProjectId, v.Version)
					find[v.ProjectId] = true
					v.Stop()
				}
			}
		}
		w.Flush()
		return nil
	},
}

func toInt(i []string) ([]int, error) {
	str := make([]int, 0, len(i))
	for _, nums := range i {
		num, err := strconv.ParseInt(nums, 10, 64)
		if err != nil {
			return nil, err
		}
		str = append(str, int(num))
	}
	return str, nil
}
