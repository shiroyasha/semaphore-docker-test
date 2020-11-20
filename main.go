package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

func ShellQuote(s string) string {
	pattern := regexp.MustCompile(`[^\w@%+=:,./-]`)

	if len(s) == 0 {
		return "''"
	}
	if pattern.MatchString(s) {
		return "'" + strings.Replace(s, "'", "'\"'\"'", -1) + "'"
	}

	return s
}

//
// code extracted from:
// https://github.com/semaphoreci/agent/blob/master/pkg/executors/docker_compose_executor.go#L265
//
func main() {
	// a random password
	password := "^&*(&#^@*&#^!@(*"

	//
	// Bug reproduction
	//
	cmd := exec.Command("bash", "-c", "echo $DOCKERHUB_PASSWORD")
	cmd.Env = []string{fmt.Sprintf("DOCKERHUB_PASSWORD=%s", ShellQuote(password))}
	bytes, _ := cmd.Output()

	// The output is '^&*(&#^@*&#^!@(*' <- notice the extra single quotes
	fmt.Println(string(bytes))

	//
	// Without ShellQuote
	//
	cmd = exec.Command("bash", "-c", "echo $DOCKERHUB_PASSWORD")
	cmd.Env = []string{fmt.Sprintf("DOCKERHUB_PASSWORD=%s", password)}
	bytes, _ = cmd.Output()

	// The output is ^&*(&#^@*&#^!@(* <- notice the lack of extra single quotes
	fmt.Println(string(bytes))

	// ------------------------------------------------------------------------
	//
	// ShellQuote is necessary for injecting environment variables into the job
	// with the `export NAME=<values>`.
	//

	//
	// However, when we pass the environment variables with cmd.Env=[]string{...}
	// the extra single quotes are unnecesary.
	//
}
