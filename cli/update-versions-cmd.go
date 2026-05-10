package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/Pumahawk/simpl-agents-controller/internal/cmd"
	"github.com/Pumahawk/simpl-agents-controller/internal/yaml"
)

var UpdateVersionCmd = cmd.Cmd{
	CName: "update-version",
	CRun: func(args []string) error {

		ref := "main"

		type prinfopath struct {
			authority    []string
			consumer     []string
			dataprovider []string
		}

		readYamlFile := func(ph string) (*yaml.Obj, error) {
			f, err := os.Open(ph)
			if err != nil {
				return nil, fmt.Errorf("unable to open file %q: %w", ph, err)
			}
			bf := &bytes.Buffer{}
			if _, err := io.Copy(bf, f); err != nil {
				return nil, err
			}
			obj := yaml.NewObj(bf.Bytes())
			return obj, nil
		}

		authorityFile := "simpl-repo/authority-iaa/charts/values.yaml"
		authorityObj, err := readYamlFile(authorityFile)
		if err != nil {
			return err
		}
		consumerFile := "simpl-repo/consumer-iaa/charts/values.yaml"
		consumerObj, err := readYamlFile(consumerFile)
		if err != nil {
			return err
		}
		dataproviderFile := "simpl-repo/provider-iaa/charts/values.yaml"
		dataproviderObj, err := readYamlFile(dataproviderFile)
		if err != nil {
			return err
		}

		var prtoupdate = map[PrInfo]prinfopath{
			prAuthenticationProvider: {
				authority: []string{"auth_provider", "targetRevision"},
			},
		}

		mapk := make([]PrInfo, 0, len(prtoupdate))
		for k := range prtoupdate {
			mapk = append(mapk, k)
		}

		find := make(map[int]bool)

		for v := range GetVersions(mapk) {
			if !find[v.PrInfo.Id] && v.Ref == ref && !regexp.MustCompile(`\.latest$`).MatchString(v.Version) {
				find[v.PrInfo.Id] = true
				v.Stop()
				path := prtoupdate[v.PrInfo]
				if path.authority != nil {
					authorityObj.UpdateAttribute(v.Version, path.authority...)
				}
				if path.consumer != nil {
					consumerObj.UpdateAttribute(v.Version, path.consumer...)
				}
				if path.dataprovider != nil {
					dataproviderObj.UpdateAttribute(v.Version, path.dataprovider...)
				}
			}
		}

		authF, err := os.OpenFile(authorityFile, os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("unable to open file authority %q: %s", authorityFile, err)
		}
		if _, err := authF.Write(authorityObj.Bytes()); err != nil {
			return fmt.Errorf("unable to update file authority %q: %s", authorityFile, err)
		}

		return nil
	},
}
