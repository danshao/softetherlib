package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// https://nathanleclaire.com/blog/2014/12/29/shelled-out-commands-in-golang/
// type SoftEtherVPNCMD struct {
// 	IP       string
// 	Password string
// }

// func (s *SoftEtherVPNCMD) Command(cmd ...string) *exec.Cmd {
// 	arg := append(
// 		[]string{
// 			fmt.Sprintf("/server %s /password:%s", s.IP, s.Password),
// 		},
// 		cmd...,
// 	)
// 	return exec.Command("vpncmd", arg...)
// }

// http://www.darrencoxall.com/golang/executing-commands-in-go/
func printCommand(cmd *exec.Cmd) {
	fmt.Printf("==> Executing: %s\n", strings.Join(cmd.Args, " "))
}

func printError(err error) {
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("==> Error: %s\n", err.Error()))
	}
}

func printOutput(outs []byte) {
	if len(outs) > 0 {
		fmt.Printf("==> Output: %s\n", string(outs))
	}
}

func main() {
	// Create an *exec.Cmd
	cmd := exec.Command("vpncmd", "/server", "127.0.0.1", "password:test", "/cmd", "serverstatus")

	// Combine stdout and stderr
	printCommand(cmd)
	output, err := cmd.CombinedOutput()
	printError(err)
	printOutput(output) // => go version go1.3 darwin/amd64
}
