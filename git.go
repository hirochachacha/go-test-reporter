package main

import (
	"bytes"
	"fmt"
	"os/exec"
)

type Git struct {
	Head        string `json:"head"`
	CommittedAt string `json:"committed_at"`
	Branch      string `json:"branch"`
}

func git() *Git {
	git := new(Git)

	cmd := exec.Command("git", "log", "-1", "--pretty=format:%H")
	head, err := cmd.Output()
	if err != nil {
		if e, ok := err.(*exec.ExitError); ok {
			panic(fmt.Sprintf("%v: %s", e, string(e.Stderr)))
		}
		panic(err)
	}
	git.Head = string(bytes.TrimSpace(head))

	cmd = exec.Command("git", "log", "-1", "--pretty=format:%ct")
	ct, err := cmd.Output()
	if err != nil {
		if e, ok := err.(*exec.ExitError); ok {
			panic(fmt.Sprintf("%v: %s", e, string(e.Stderr)))
		}
		panic(err)
	}
	git.CommittedAt = string(bytes.TrimSpace(ct))

	cmd = exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	br, err := cmd.Output()
	if err != nil {
		if e, ok := err.(*exec.ExitError); ok {
			panic(fmt.Sprintf("%v: %s", e, string(e.Stderr)))
		}
		panic(err)
	}
	git.Branch = string(bytes.TrimSpace(br))

	return git
}
