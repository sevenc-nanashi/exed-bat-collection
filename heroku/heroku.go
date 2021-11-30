// +build ignore
// path: .../heroku/bin/heroku.exe

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
	os.Setenv("HEROKU_BINPATH", path.Join(dp0, "heroku.cmd"))
	var app string
	var args []string
	if Exists(os.Getenv("LOCALAPPDATA") + "/heroku/client/bin/heroku.exe") {
		app = os.Getenv("LOCALAPPDATA") + "/heroku/client/bin/heroku.exe"
		args = os.Args[1:]
	} else {
		app = dp0 + "../client/bin/node.exe"
		args = append([]string{dp0 + "../client/bin/run"}, os.Args[1:]...)
	}
	trap := make(chan os.Signal, 1)
	signal.Ignore(syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT)
	main_cmd := exec.Command(app, args...)
	main_cmd.Stdin = os.Stdin
	main_cmd.Stdout = os.Stdout
	main_cmd.Stderr = os.Stderr
	main_cmd.Run()
	os.Exit(main_cmd.ProcessState.ExitCode())

}
