package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
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
		&LastVersionCmd,
		&UpdateVersionCmd,
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

		if len(fs.Args()) == 0 {
			return fmt.Errorf("missing project id")
		}

		id, err := getId(fs.Arg(0))
		if err != nil {
			return err
		}

		items, err := glclient.Packages(int(id), page, perpage)
		if err != nil {
			return err
		}

		w := newTabWriter()
		fmt.Fprintf(w, "Id\tName\tVersion\tType\tPipelineId\tRef\tPipeline WebUrl\n")
		for i := range items {
			item := items[i]
			fmt.Fprintf(w, "%d\t%q\t%q\t%q\t%d\t%q\t%q\n",
				item.Id, item.Name, item.Version, item.PackageType, item.Pipeline.Id, item.Pipeline.Ref, item.Pipeline.WebUrl)
		}
		w.Flush()
		return nil
	},
}

var LastVersionCmd = cmd.Cmd{
	CName: "last-version",
	CRun: func(args []string) error {
		var ref string
		var num int
		fs := flag.NewFlagSet("", flag.ExitOnError)
		fs.IntVar(&num, "num", 1, "")
		fs.StringVar(&ref, "ref", "main", "")
		fs.Parse(args)

		projectIds := prIdsDemux.Demux(fs.Args())

		if len(projectIds) == 0 {
			return fmt.Errorf("missing project id")
		}

		nums, err := getPrsInfo(projectIds)
		if err != nil {
			return err
		}

		cv := GetVersions(nums)
		w := newTabWriter()
		find := make(map[int]int)
		fmt.Fprintf(w, "Project\tName\tRef\tVersion\n")
		for v := range cv {
			if v.Ref == ref && v.Type == v.PrInfo.Type && !regexp.MustCompile(`\.latest$`).MatchString(v.Version) {
				if find[v.PrInfo.Id] < num {
					fmt.Fprintf(w, "%d\t%q\t%q\t%q\n", v.PrInfo.Id, v.PrInfo.Name, v.Ref, v.Version)
					find[v.PrInfo.Id]++
				}
				if find[v.PrInfo.Id] >= num {
					v.Stop()

				}
			}
		}
		w.Flush()
		return nil
	},
}

func getId(idx string) (int, error) {
	var id int
	if ids, ok := prIds.Get(idx); ok {
		id = ids.Id
	} else {
		idx, err := strconv.ParseInt(idx, 10, 64)
		if err != nil {
			return 0, err
		}
		id = int(idx)
	}
	return id, nil
}

func getPrsInfo(i []string) ([]PrInfo, error) {
	str := make([]PrInfo, 0, len(i))
	for _, nums := range i {
		info, ok := prIds.Get(nums)
		if !ok {
			return nil, fmt.Errorf("unknown project %q", i)
		}
		str = append(str, info)
	}
	return str, nil
}

func newTabWriter() *tabwriter.Writer {
	return tabwriter.NewWriter(os.Stdout, 1, 2, 3, ' ', 0)
}
