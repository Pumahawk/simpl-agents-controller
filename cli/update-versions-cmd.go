package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"slices"
	"strings"

	"github.com/Pumahawk/simpl-agents-controller/internal/cmd"
	"github.com/Pumahawk/simpl-agents-controller/internal/yaml"
)

var UpdateVersionCmd = cmd.Cmd{
	CName: "update-version",
	CRun: func(args []string) error {

		var ref string

		fc := flag.NewFlagSet("", flag.ExitOnError)
		fc.StringVar(&ref, "ref", "main", "")
		fc.Parse(args)

		readYamlFile := func(ph string) (*yaml.Obj, error) {
			f, err := os.Open(ph)
			if err != nil {
				return nil, fmt.Errorf("unable to open file %q: %w", ph, err)
			}
			defer f.Close()
			bf := &bytes.Buffer{}
			if _, err := io.Copy(bf, f); err != nil {
				return nil, err
			}
			obj := yaml.NewObj(bf.Bytes())
			return obj, nil
		}

		conf := map[string]map[PrInfo][]string{
			"simpl-repo/authority-iaa/charts/values.yaml": {
				prAuthenticationProvider:      {"auth_provider", "targetRevision"},
				prEidasKeycloak:               {"keycloak", "eidas_demo_keycloak_extension", "targetRevision"},
				prFeAuthenticationProvider:    {"fe_auth_provider", "targetRevision"},
				prFeIdentityProvider:          {"fe_identity_provider", "targetRevision"},
				prFeOnboarding:                {"fe_onboarding", "targetRevision"},
				prFeSecurityAttributeProvider: {"fe_security_attribute_provider", "targetRevision"},
				prFeUsersAndRoles:             {"fe_users_roles", "targetRevision"},
				prIdentityProvider:            {"identity_provider", "targetRevision"},
				prKeycloakAuthenticator:       {"keycloak", "keycloak_authenticator", "targetRevision"},
				prOnboarding:                  {"onboarding", "targetRevision"},
				prSecurityAttributesProvider:  {"sap", "targetRevision"},
				prTier1Gateway:                {"tier1_gateway", "targetRevision"},
				prTier2Gateway:                {"tier2_gateway", "targetRevision"},
				prTier2Proxy:                  {"tier2_proxy", "targetRevision"},
				prUsersRoles:                  {"users_roles", "targetRevision"},
			},
			"simpl-repo/consumer-iaa/charts/values.yaml": {
				prAuthenticationProvider:   {"auth_provider", "targetRevision"},
				prFeAuthenticationProvider: {"auth_provider_fe", "targetRevision"},
				prFeAuthenticationProvider: {"authentication_provider_fe", "targetRevision"},
				prFeUsersAndRoles:          {"users_roles_fe", "targetRevision"},
				prKeycloakAuthenticator:    {"keycloak", "keycloak_authenticator", "targetRevision"},
				prTier1Gateway:             {"tier1_gateway", "targetRevision"},
				prTier2Gateway:             {"tier2_gateway", "targetRevision"},
				prTier2Proxy:               {"tier2_proxy", "targetRevision"},
				prUsersRoles:               {"users_roles", "targetRevision"},
			},
			"simpl-repo/provider-iaa/charts/values.yaml": {
				prAuthenticationProvider:   {"auth_provider", "targetRevision"},
				prFeAuthenticationProvider: {"auth_provider_fe", "targetRevision"},
				prFeAuthenticationProvider: {"authentication_provider_fe", "targetRevision"},
				prFeUsersAndRoles:          {"users_roles_fe", "targetRevision"},
				prKeycloakAuthenticator:    {"keycloak", "keycloak_authenticator", "targetRevision"},
				prTier1Gateway:             {"tier1_gateway", "targetRevision"},
				prTier2Gateway:             {"tier2_gateway", "targetRevision"},
				prTier2Proxy:               {"tier2_proxy", "targetRevision"},
				prUsersRoles:               {"users_roles", "targetRevision"},
			},
		}

		type yamlSt struct {
			file string
			obj  *yaml.Obj
		}
		yamlObjs := make(map[string]*yamlSt)
		for f := range conf {
			obj, err := readYamlFile(f)
			if err != nil {
				return err
			}
			yamlObjs[f] = &yamlSt{f, obj}
		}

		yamlObjFiles := make(map[PrInfo][]string)
		for f, v := range conf {
			for info := range v {
				yamlObjFiles[info] = append(yamlObjFiles[info], f)
			}
		}

		mapk := make([]PrInfo, 0)
		for _, v := range conf {
			for info := range v {
				if !slices.Contains(mapk, info) {
					mapk = append(mapk, info)
				}
			}
		}

		find := make(map[int]bool)

		for v := range GetVersions(mapk) {
			if !find[v.PrInfo.Id] && v.Ref == ref && !regexp.MustCompile(`\.latest$`).MatchString(v.Version) {
				find[v.PrInfo.Id] = true
				v.Stop()
				for _, fl := range yamlObjFiles[v.PrInfo]{
				obj, ok := yamlObjs[fl]
				if ok {
					path, ok := conf[obj.file][v.PrInfo]
					if ok {
						if ok, err := obj.obj.UpdateAttribute(v.Version, path...); err != nil {
							return fmt.Errorf("unable to update file=%q path=%q", obj.file, strings.Join(path, "."))
						} else if !ok {
							fmt.Fprintf(os.Stderr, "unable to update file=%q %q version %q\n", obj.file, v.PrInfo.Name, v.Version)
						}
					}
				}
				}
			}
		}

		for file := range conf {
			fopened, err := os.OpenFile(file, os.O_TRUNC|os.O_WRONLY, 0644)
			if err != nil {
				return fmt.Errorf("unable to open file %q: %s", file, err)
			}
			defer fopened.Close()
			for _, obj := range yamlObjs {
				if obj.file == file {
					if _, err := fopened.Write(obj.obj.Bytes()); err != nil {
						return fmt.Errorf("unable to update file=%q: %s", file, err)
					}
				}
			}
		}

		return nil
	},
}
