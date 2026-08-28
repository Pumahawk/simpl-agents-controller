package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"maps"
	"os"
	"os/exec"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"sync"

	"github.com/Pumahawk/simpl-agents-controller/internal/cmd"
)

var RemoteUpdateCmd = cmd.Cmd{
	CName: "remote:update",
	CRun: func(args []string) error {
		var apply, syncv bool
		var jobs int

		fs := flag.NewFlagSet("", flag.ExitOnError)
		fs.BoolVar(&apply, "apply", false, "")
		fs.BoolVar(&syncv, "sync", false, "")
		fs.IntVar(&jobs, "jobs", 5, "")
		fs.Parse(args)
		args = fs.Args()

		depInfo := parseArgsDeployer(args)
		deployers := slices.Collect(maps.Keys(depInfo))
		argoApp, err := GetArgoApps(deployers)
		if err != nil {
			fmt.Printf("errore on get argo apps: %s\n", err)
			os.Exit(1)
		}
		ps := argoApp.GetParameters()
		ids := argoApp.GetProjectIds()
		versionsStore := make(map[int]map[string]string)
		for _, p := range ps {
			if _, ok := depInfo[p.Deployer]; ok {
				if pr, ok := projectIdByTargetName[p.Name]; ok {
					mf, ok := versionsStore[pr.ProjectId]
					if !ok {
						mf = make(map[string]string)
					}
					br := depInfo[p.Deployer]
					mf[br] = ""
					versionsStore[pr.ProjectId] = mf
				}
			}
		}
		for d, id := range ids {
			mf, ok := versionsStore[id]
			if !ok {
				mf = make(map[string]string)
			}
			br := depInfo[d]
			mf[br] = ""
			versionsStore[id] = mf
		}
		GetProjectLastVersion(versionsStore)
		argoApp.UpdateParameters(versionsStore, depInfo)
		if apply {
			c := exec.Command("kubectl", "apply", "-f", "-")
			c.Stderr = os.Stderr
			c.Stdout = os.Stdout
			out, err := json.Marshal(argoApp)
			if err != nil {
				fmt.Printf("unable to create json from argoApp: %s\n", err)
				os.Exit(1)
			}
			c.Stdin = bytes.NewBuffer(out)
			if err := c.Run(); err != nil {
				os.Exit(1)
			}
		} else {
			json.NewEncoder(os.Stdout).Encode(argoApp)
		}

		if syncv {
			poolsize := make(chan int, jobs)
			wg := &sync.WaitGroup{}
			for i := range deployers {
				deployer := deployers[i]
				wg.Go(func() {
					poolsize <- 0
					defer func() {
						<-poolsize
					}()
					c := exec.Command("kubectl", "-n", "argocd", "exec", "argo-cd-argocd-application-controller-0", "--", "argocd", "--core", "app", "sync", "--prune", "--async", deployer)
					serr := &bytes.Buffer{}
					c.Stderr = serr
					if err := c.Run(); err != nil {
						fmt.Printf("synch error %q: %s\n", deployer, serr)
					} else {
						fmt.Printf("synched %q\n", deployer)
					}
				})
			}
			wg.Wait()
		}
		return nil
	},
}

type ArgoApps map[string]any

func GetArgoApps(apps []string) (ArgoApps, error) {
	out := &bytes.Buffer{}
	args := append([]string{"-n", "argocd", "get", "application", "-o", "json"}, apps...)
	c := exec.Command("kubectl", args...)
	c.Stderr = os.Stderr
	c.Stdout = out
	if err := c.Run(); err != nil {
		return nil, err
	}
	data := make(ArgoApps)
	if err := json.NewDecoder(out).Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}

type Paramater struct {
	Deployer string
	Name     string
	Value    string
}

func (a ArgoApps) GetProjectIds() map[string]int {
	out := make(map[string]int)
	items := []any{map[string]any(a)}
	if v, ok := a["items"]; ok && v != nil {
		if v, ok := v.([]any); ok {
			items = v
		}
	}
	for _, item := range items {
		if item, ok := item.(map[string]any); ok {
			if metadata, ok := item["metadata"]; ok && metadata != nil {
				if metadata, ok := metadata.(map[string]any); ok {
					if deployer, ok := metadata["name"]; ok && deployer != nil {
						if deployer, ok := deployer.(string); ok {
							if v, ok := item["spec"]; ok && v != nil {
								if v, ok := v.(map[string]any); ok {
									if v, ok := v["source"]; ok && v != nil {
										if v, ok := v.(map[string]any); ok {
											if repoURL, ok := v["repoURL"]; ok {
												if repoURL, ok := repoURL.(string); ok {
													if id, err := projectIdFromRepoUrl(repoURL); err == nil {
														out[deployer] = id
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return out
}

func (a ArgoApps) GetParameters() []Paramater {
	var p []Paramater
	items := []any{map[string]any(a)}
	if v, ok := a["items"]; ok && v != nil {
		if v, ok := v.([]any); ok {
			items = v
		}
	}
	for _, item := range items {
		if item, ok := item.(map[string]any); ok {
			if metadata, ok := item["metadata"]; ok && metadata != nil {
				if metadata, ok := metadata.(map[string]any); ok {
					if deployer, ok := metadata["name"]; ok && deployer != nil {
						if deployer, ok := deployer.(string); ok {
							if v, ok := item["spec"]; ok && v != nil {
								if v, ok := v.(map[string]any); ok {
									if v, ok := v["source"]; ok && v != nil {
										if v, ok := v.(map[string]any); ok {
											if v, ok := v["helm"]; ok && v != nil {
												if v, ok := v.(map[string]any); ok {
													if v, ok := v["parameters"]; ok && v != nil {
														if v, ok := v.([]any); ok {
															for _, parameter := range v {
																if parameter, ok := parameter.(map[string]any); ok {
																	if name, ok := parameter["name"]; ok && name != nil {
																		if name, ok := name.(string); ok {
																			if value, ok := parameter["value"]; ok && value != nil {
																				if value, ok := value.(string); ok {
																					p = append(p, Paramater{deployer, name, string(value)})
																				}
																			}
																		}
																	}
																}
															}
														}
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return p
}

func (a ArgoApps) UpdateParameters(versionsStore map[int]map[string]string, depInfo map[string]string) {
	items := []any{map[string]any(a)}
	if v, ok := a["items"]; ok && v != nil {
		if v, ok := v.([]any); ok {
			items = v
		}
	}
	for _, item := range items {
		if item, ok := item.(map[string]any); ok {
			if metadata, ok := item["metadata"]; ok && metadata != nil {
				if metadata, ok := metadata.(map[string]any); ok {
					if deployer, ok := metadata["name"]; ok && deployer != nil {
						if deployer, ok := deployer.(string); ok {
							if v, ok := item["spec"]; ok && v != nil {
								if v, ok := v.(map[string]any); ok {
									if v, ok := v["source"]; ok && v != nil {
										if v, ok := v.(map[string]any); ok {
											if repoURL, ok := v["repoURL"]; ok {
												if repoURL, ok := repoURL.(string); ok {
													if id, err := projectIdFromRepoUrl(repoURL); err == nil {
														branch := depInfo[deployer]
														version := versionsStore[id][branch]
														if version != "" {
															v["targetRevision"] = version
														}
													}
												}
											}
											if v, ok := v["helm"]; ok && v != nil {
												if v, ok := v.(map[string]any); ok {
													if v, ok := v["parameters"]; ok && v != nil {
														if v, ok := v.([]any); ok {
															for _, parameter := range v {
																if parameter, ok := parameter.(map[string]any); ok {
																	if name, ok := parameter["name"]; ok && name != nil {
																		if name, ok := name.(string); ok {
																			if value, ok := parameter["value"]; ok && value != nil {
																				if _, ok := value.(string); ok {
																					if info, ok := projectIdByTargetName[name]; ok {
																						branch := depInfo[deployer]
																						version := versionsStore[info.ProjectId][branch]
																						if version != "" {
																							parameter["value"] = version
																						}
																					}
																				}
																			}
																		}
																	}
																}
															}
														}
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
}

type ProjectInfo struct {
	Type      string
	ProjectId int
}

var projectIdByTargetName = map[string]ProjectInfo{
	"auth_provider.targetRevision":                          {"helm", 939},
	"authentication_provider_fe.targetRevision":             {"helm", 1308},
	"consent_management.targetRevision":                     {"helm", 1723},
	"consent_management_fe.targetRevision":                  {"helm", 1950},
	"fe_auth_provider.targetRevision":                       {"helm", 1308},
	"fe_identity_provider.targetRevision":                   {"helm", 1311},
	"fe_onboarding.targetRevision":                          {"helm", 1307},
	"fe_security_attribute_provider.targetRevision":         {"helm", 1309},
	"fe_users_roles.targetRevision":                         {"helm", 1999},
	"identity_provider.targetRevision":                      {"helm", 2119},
	"keycloak.eidas_demo_keycloak_extension.targetRevision": {"maven", 1313},
	"keycloak.keycloak_authenticator.targetRevision":        {"maven", 915},
	"keycloak.oid4vp_keycloak_extension.targetRevision":     {"maven", 1840},
	"onboarding.targetRevision":                             {"helm", 2097},
	"sap.targetRevision":                                    {"helm", 861},
	"tier1_gateway.targetRevision":                          {"helm", 2112},
	"tier2_gateway.targetRevision":                          {"helm", 2215},
	"tier2_proxy.targetRevision":                            {"helm", 1112},
	"users_roles.targetRevision":                            {"helm", 2000},
	"users_roles_fe.targetRevision":                         {"helm", 1999},
}

// map[projectId]map[ref]version
func GetProjectLastVersion(out map[int]map[string]string) {
	var pris []PrInfo
	for prId := range out {
		pris = append(pris, PrInfo{prId, "nd", "nd"})
	}
	chv := GetVersions(pris)

main:
	for pk := range chv {
		if !strings.HasSuffix(pk.Version, ".latest") {
			if v, ok := out[pk.PrInfo.Id][pk.Ref]; ok {
				if v == "" {
					out[pk.PrInfo.Id][pk.Ref] = pk.Version
					for _, v := range out[pk.PrInfo.Id] {
						if v == "" {
							continue main
						}
					}
					pk.Stop()
				}
			}
		}
	}
}

func parseArgsDeployer(args []string) map[string]string {
	out := make(map[string]string)
	for _, arg := range args {
		v := strings.Split(arg, ":")
		deployer := v[0]
		branch := "main"
		if len(v) > 1 {
			branch = v[1]
		}
		out[deployer] = branch
	}
	return out
}

func projectIdFromRepoUrl(repoURL string) (int, error) {
	rx := regexp.MustCompile("https://code.europa.eu/api/v4/projects/([0-9][0-9]*)/packages/helm/stable")
	m := rx.FindStringSubmatch(repoURL)
	if len(m) == 2 {
		return strconv.Atoi(m[1])
	}
	return -1, fmt.Errorf("unable to retrieve project id from repoURL %q", repoURL)
}
