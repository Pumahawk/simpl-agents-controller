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
	CName: "bao:get",
	CRun: func(args []string) error {
		var version int
		fs := flag.NewFlagSet("", flag.ExitOnError)
		fs.IntVar(&version, "v", -1, "")
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

		if fs.NArg() > 0 {
			if fs.NArg() == 1 {
				KeyList(cl, fs.Arg(0))
			} else {
				Secret(cl, fs.Arg(0), fs.Arg(1), version)
			}
		} else {
			MountCmd(cl)
		}
		return nil
	},
}

func MountCmd(cl *bao.Client) {
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
}

func KeyList(cl *bao.Client, key string) {
	res, err := cl.KeysList(key)
	if err != nil {
		fmt.Println("unable to get list keys", err)
		os.Exit(1)
	}
	for _, v := range res.Items {
		fmt.Println(v.Name)
	}
}

func Secret(cl *bao.Client, key, name string, version int) {
	resmeta, err := cl.SecretVers(key, name)
	if err != nil {
		fmt.Printf("unable to get secret key=%q, name=%q: %s\n", key, name, err)
		os.Exit(1)
	}
	if version == -1 {
		version = resmeta.CurrentVersion
	}
	res, err := cl.SecretVer(key, name, version)
	if err != nil {
		fmt.Printf("unable to get secret with version key=%q, name=%q vers=%d: %s\n", key, name, version)
		os.Exit(1)
	}

	tw := tabwriter.NewWriter(os.Stdout, 2, 2, 2, ' ', 0)
	fmt.Fprintf(tw, "%s\t%d\n", "Current Version:", version)
	fmt.Fprintf(tw, "%s\t%d\n", "Oldest Version:", resmeta.OldestVersion)
	tw.Flush()
	fmt.Fprintf(tw, "%s\t%s\n", "Name", "Secret")
	for _, v := range res.Items {
		fmt.Fprintf(tw, "%s\t%s\n", v.Key, v.Value)
	}
	tw.Flush()
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
