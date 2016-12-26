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
	head, err := cmd.CombinedOutput()
	if err != nil {
		panic(fmt.Sprintf("%v: %v", err, string(head)))
	}
	git.Head = string(bytes.TrimSpace(head))

	cmd = exec.Command("git", "log", "-1", "--pretty=format:%ct")
	ct, err := cmd.Output()
	if err != nil {
		panic(fmt.Sprintf("%v: %v", err, string(ct)))
	}
	git.CommittedAt = string(bytes.TrimSpace(ct))

	cmd = exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	br, err := cmd.CombinedOutput()
	if err != nil {
		panic(fmt.Sprintf("%v: %v", err, string(br)))
	}
	git.Branch = string(bytes.TrimSpace(br))

	return git
}
