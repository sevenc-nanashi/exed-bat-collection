// +build ignore
// path: %LOCALAPPDATA%/heroku/client/bin/heroku.exe

package main

import (
	"io/ioutil"
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
	trap := make(chan os.Signal, 1)
	signal.Notify(trap, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT)
	go func() {
		<-trap
		// Do nothing
	}()
	os.Setenv("HEROKU_REDIRECTED", "1")
	os.Setenv("HEROKU_BINPATH", dp0+"/heroku")
	cmd_content_bytes, _ := ioutil.ReadFile(dp0 + "/heroku.cmd")
	cmd_content := string(cmd_content_bytes)
	raw_path := strings.Split(cmd_content, "\"")[1]
	app := strings.Replace(raw_path, "\\", "/", -1)
	app = strings.Replace(app, "%~dp0", dp0+"/", -1)
	app = strings.Replace(app, "heroku.cmd", "heroku.exe", -1)
	args := os.Args[1:]
	main_cmd := exec.Command(app, args...)
	main_cmd.Stdin = os.Stdin
	main_cmd.Stdout = os.Stdout
	main_cmd.Stderr = os.Stderr
	main_cmd.Run()
	os.Exit(main_cmd.ProcessState.ExitCode())
}
