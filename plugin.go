package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type (
	Plugin struct {
		Token   string
		Project string
		Message string
		Targets string
		DryRun  bool
		Debug   bool
	}
)

func (p Plugin) Exec() error {
	if p.shouldSetProject() {
		use := p.buildUse()
		if err := execute(use, p.Debug, p.DryRun); err != nil {
			return err
		}
	}

	deploy := p.buildDeploy()
	if err := execute(deploy, p.Debug, p.DryRun); err != nil {
		return err
	}

	return nil
}

func (p *Plugin) shouldSetProject() bool {
	return p.Project != ""
}

func getEnvironment(oldEnv []string, p *Plugin) []string {
	var env []string
	for _, v := range oldEnv {
		if !strings.HasPrefix(v, "DEBUG=") && !strings.HasPrefix(v, "FIREBASE_TOKEN=") {
			env = append(env, v)
		}
	}
	env = append(env, fmt.Sprintf("FIREBASE_TOKEN=%s", p.Token))
	if p.Debug {
		env = append(env, fmt.Sprintf("DEBUG=%s", "true"))
	}
	return env
}

// buildUse creates a command on the form:
// $ firebase use ...
func (p *Plugin) buildUse() *exec.Cmd {
	var args []string
	args = append(args, "use")

	if p.Project != "" {
		args = append(args, p.Project)
	}

	cmd := exec.Command("firebase", args...)
	cmd.Env = getEnvironment(os.Environ(), p)
	return cmd
}

// buildDeploy creates a command on the form:
// $ firebase deploy \
//   [--only ...] \
//   [--message ...]
func (p *Plugin) buildDeploy() *exec.Cmd {
	var args []string
	args = append(args, "deploy")

	if p.Targets != "" {
		args = append(args, "--only")
		args = append(args, p.Targets)
	}

	if p.Message != "" {
		args = append(args, "--message")
		args = append(args, fmt.Sprintf("\"%s\"", p.Message))
	}

	cmd := exec.Command("firebase", args...)
	cmd.Env = getEnvironment(os.Environ(), p)
	return cmd
}

// execute sets the stdout and stderr of the command to be the default, traces
// the command to be executed and returns the result of the command execution.
func execute(cmd *exec.Cmd, debug, dryRun bool) error {
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if dryRun || debug {
		fmt.Println("$", strings.Join(cmd.Args, " "))
	}
	if dryRun {
		return nil
	}
	return cmd.Run()
}
