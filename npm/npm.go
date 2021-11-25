// +build ignore
// path: .../nodejs/npm.exe

package main

import (
	"os"
	"os/exec"
	"os/signal"
	"path"
	"strings"
	"syscall"
)

func Exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func main() {
	executable, _ := os.Executable()
	executable = strings.Replace(executable, "\\", "/", -1)
	dp0 := path.Dir(executable)
	var node_exe string
	var args []string
	if Exists(dp0 + "/node.exe") {
		node_exe = dp0 + "/node.exe"
	} else {
		node_exe = "node"
	}
	npm_cli_js := dp0 + "/node_modules/npm/bin/npm-cli.js"
	npm_prefix_npm_cli_js_cmd := exec.Command(node_exe, npm_cli_js, "prefix", "-g")
	npm_prefix_npm_cli_js_out, _ := npm_prefix_npm_cli_js_cmd.Output()
	npm_prefix_npm_cli_js := string(npm_prefix_npm_cli_js_out) + "/node_modules/npm/bin/npm-cli.js"
	if Exists(npm_prefix_npm_cli_js) {
		npm_cli_js = npm_prefix_npm_cli_js
	}
	args = append(args, npm_cli_js)
	args = append(args, os.Args[1:]...)
	trap := make(chan os.Signal, 1)
	signal.Notify(trap, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT)
	go func() {
		<-trap
		// Do nothing
	}()
	main_cmd := exec.Command(node_exe, args...)
	main_cmd.Stdin = os.Stdin
	main_cmd.Stdout = os.Stdout
	main_cmd.Stderr = os.Stderr
	main_cmd.Run()
	os.Exit(main_cmd.ProcessState.ExitCode())
}
