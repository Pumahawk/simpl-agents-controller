package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// Flags
var skipJavaTests = false

// Dependencies
var depJava = "java@21"
var depMaven = "mvnd@1.0"
var depHelm = "helm@4.1.1"
var depDocker = "docker-cli@28.4"

var dockerImage string

func main() {
	if err := LoadFlags(); err != nil {
		log.Fatalf("Invalid flags: %s", err)
	}
	if err := Doctor(); err != nil {
		log.Fatalf("Invalid environment: %s", err)
	}

	command := flag.Arg(0)
	if command == "" {
		log.Fatalf("Missing command")
	}
	switch command {
	case "java-build":
		if err := JavaBuild(); err != nil {
			log.Fatalf("Java compilation error: %s", err)
		}
	case "helm-build":
		if err := HelmBuild(); err != nil {
			log.Fatalf("Helm compilation error: %s", err)
		}
	case "docker-build":
		if err := DockerBuild(); err != nil {
			log.Fatalf("Docker compilation error: %s", err)
		}
	default:
		log.Fatalf("Unknown command %s", command)
	}
}

func Doctor() error {
	if err := exec.Command("mise", "--version").Run(); err != nil {
		return fmt.Errorf("missing mise dependency: %w", err)
	}
	if err := exec.Command("git", "--version").Run(); err != nil {
		return fmt.Errorf("missing git dependency: %w", err)
	}
	return nil
}

func LoadFlags() error {
	flag.BoolVar(&skipJavaTests, "skip-tests", false, "")
	flag.Parse()
	return nil
}

func JavaBuild() error {
	miseargs := []string{
		"x",
		depJava, depMaven,
		"--",
	}
	javaBuildArgs := []string{"mvnd", "clean", "install"}
	if skipJavaTests {
		javaBuildArgs = append(javaBuildArgs, "-Dmaven.test.skip")
	}
	args := append(miseargs, javaBuildArgs...)
	cmd := exec.Command("mise",
		args...,
	)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("java build error: %w", err)
	}
	return nil
}

func HelmBuild() error {
	miseargs := []string{
		"x",
		depHelm,
		"--",
	}
	helmBuildArgs := []string{"helm", "package", "charts"}
	args := append(miseargs, helmBuildArgs...)
	cmd := exec.Command("mise",
		args...,
	)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("helm build error: %w", err)
	}
	return nil
}

func DockerBuild() error {
	info, err := GetRepoInfo()
	if err != nil {
		return fmt.Errorf("unable to retrieve project info: %w", err)
	}
	vers, err := GetProjectVersionFromPipelineVariablesFile()
	if err != nil {
		return fmt.Errorf("unable to retrieve project version from pipeline file: %w", err)
	}

	dockerImage := RetrieveDockerImage(info, vers)

	miseargs := []string{
		"x",
		depDocker,
		"--",
	}
	dockerBuildArgs := []string{"docker", "build", "-t", dockerImage, "."}
	args := append(miseargs, dockerBuildArgs...)
	cmd := exec.Command("mise",
		args...,
	)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("docker build error: %w", err)
	}
	return nil
}

func GetProjectIdFromGitRemote() (string, error) {
	out := &bytes.Buffer{}
	cmd := exec.Command("git", "remote", "get-url", "origin")
	cmd.Stdout = out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("unable to execute git remote get-url origin: %w", err)
	}
	ru := out.String()
	ru = strings.Trim(ru, "\n\r \t")
	if !strings.HasPrefix(ru, "https://code.europa.eu/") {
		return "", fmt.Errorf("unknown git repository host %q: %w")
	}
	ru = regexp.MustCompile("^https://code.europa.eu/").ReplaceAllString(ru, "")
	ru = regexp.MustCompile("\\.git/?$").ReplaceAllString(ru, "")
	return url.QueryEscape(ru), nil
}

type RepoInfo struct {
	Id                int    `json:"id"`
	Name              string `json:"name"`
	Path              string `json:"path"`
	PathWithNamespace string `json:"path_with_namespace"`
	WebUrl            string `json:"web_url"`
}

func GetRepoInfo() (*RepoInfo, error) {
	cl := &http.Client{}
	prid, err := GetProjectIdFromGitRemote()
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf("https://code.europa.eu/api/v4/projects/%s", prid)
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	res, err := cl.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 && res.StatusCode >= 300 {
		return nil, fmt.Errorf("unable to get project path code=%d projectID=%q: %w", res.StatusCode, prid, err)
	}
	rinfo := &RepoInfo{}
	if err := json.NewDecoder(res.Body).Decode(rinfo); err != nil {
		return nil, err
	}
	return rinfo, nil
}

func GetProjectVersionFromPipelineVariablesFile() (string, error) {
	f, err := os.Open("pipeline.variables.sh")
	if err != nil {
		return "", err
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		if strings.Contains(line, "PROJECT_VERSION_NUMBER") {
			fd := regexp.MustCompile("=\"(.*)\"$").FindStringSubmatch(line)
			if len(fd) > 1 {
				return fd[1], nil
			} else {
				return "", fmt.Errorf("unable to extract version from PROJECT_VERSION_NUMBER line")
			}
		}
	}
	return "", fmt.Errorf("not found version in pipeline attributes file")
}

func RetrieveDockerImage(info *RepoInfo, version string) string {
	if dockerImage == "" {
		return GenerateDockerImage(info, version)
	} else {
		return dockerImage
	}
}

func GenerateDockerImage(info *RepoInfo, version string) string {
	return fmt.Sprintf("code.europa.eu/%s:%s-local.0.00000000", info.PathWithNamespace, version)
}
