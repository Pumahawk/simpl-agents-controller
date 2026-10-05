package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
)

// Flags
var skipJavaTests = false

// Dependencies
var depJava = "java@21"
var depMaven = "mvnd@1.0"
var depHelm = "helm@4.1.1"
var depDocker = "docker-cli@28.4"

var dockerImage string = "temp-image-docker"

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
