package main

import (
	"flag"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/Pumahawk/simpl-agents-controller/internal/bao"
	"github.com/Pumahawk/simpl-agents-controller/internal/cmd"
)

var TokenBaoCmd = cmd.Cmd{
	CName: "bao:token",
	CRun: func(args []string) error {

		fs := flag.NewFlagSet("", flag.ExitOnError)
		cfg := BaoGFlags(fs)
		fs.Parse(args)

		cl := cfg.KToken()

		tk, err := cl.GetToken()
		if err != nil {
			fmt.Println("bao token error", err)
			os.Exit(1)
		}
		fmt.Print(tk)
		return nil
	},
}

var ListBaoCmd = cmd.Cmd{
	CName: "bao:list",
	CRun: func(args []string) error {

		fs := flag.NewFlagSet("", flag.ExitOnError)
		cfg := BaoGFlags(fs)
		fs.Parse(args)

		clk := cfg.KToken()

		tk, err := clk.GetToken()
		if err != nil {
			fmt.Println("bao token error", err)
			os.Exit(1)
		}

		cl := cfg.Client()
		cl.TokenFunc = func() (string, error) {
			return tk, nil
		}

		if cl.Url == "" {
			host, err := clk.GetSecretHost()
			if err != nil {
				fmt.Println("get secret host from kube", err)
				os.Exit(1)
			}
			cl.Url = host
		}

		res, err := cl.Mounts()
		if err != nil {
			fmt.Println("bao mounts error", err)
			os.Exit(1)
		}

		tw := tabwriter.NewWriter(os.Stdout, 2, 2, 2, ' ', 0)
		fmt.Fprintf(tw, "%s\t%s\n", "Name", "Desc")
		for _, v := range res.Items {
			fmt.Fprintf(tw, "%s\t%s\n", v.Name, v.Desc)
		}
		tw.Flush()

		return nil
	},
}

type BaoFlags struct {
	Url       string
	Namespace string
}

func (b *BaoFlags) KToken() *bao.KClient {
	cl := &bao.KClient{}
	cl.Namespace = b.Namespace
	return cl
}

func (b *BaoFlags) Client() *bao.Client {
	cl := &bao.Client{}
	cl.Url = b.Url
	return cl
}

func BaoGFlags(fs *flag.FlagSet) *BaoFlags {
	bf := &BaoFlags{}
	fs.StringVar(&bf.Namespace, "ns", "common01", "")
	fs.StringVar(&bf.Url, "url", "", "")
	return bf
}
