package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/Pumahawk/simpl-agents-controller/internal/cmd"
	"github.com/Pumahawk/simpl-agents-controller/internal/keycloak"
)

var TokenizeCmd = cmd.Cmd{
	CName: "keycloak:tokenize",
	CRun: func(args []string) error {
		cl := keycloak.Client{}
		var ojson bool
		fc := flag.NewFlagSet("", flag.ExitOnError)
		fc.StringVar(&cl.Server, "server", "", "")
		fc.StringVar(&cl.Realm, "realm", "", "")
		fc.StringVar(&cl.ClientId, "clientid", "", "")
		fc.StringVar(&cl.ClientSecret, "clientsecret", "", "")
		fc.StringVar(&cl.User, "user", "", "")
		fc.StringVar(&cl.Password, "password", "", "")
		fc.BoolVar(&ojson, "json", false, "")
		fc.Parse(args)

		res, err := cl.Tokenize()
		if err != nil {
			fmt.Println("tokenize error", err)
			os.Exit(1)
		}

		if ojson {
			jenc := json.NewEncoder(os.Stdout)
			jenc.SetIndent("", "  ")
			jenc.Encode(res)
		} else {
			fmt.Printf("%s", res.AccessToken)
		}
		return nil
	},
}
