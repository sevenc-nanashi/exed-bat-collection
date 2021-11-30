// +build ignore
// path: %LOCALAPPDATA%/heroku/client/(version)/heroku.exe

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

func startExe(app string, args []string) {
	main_cmd := exec.Command(app, args...)
	main_cmd.Stdin = os.Stdin
	main_cmd.Stdout = os.Stdout
	main_cmd.Stderr = os.Stderr
	main_cmd.Run()
	os.Exit(main_cmd.ProcessState.ExitCode())
}

func main() {
	executable, _ := os.Executable()
	executable = strings.Replace(executable, "\\", "/", -1)
	dp0 := path.Dir(executable)
	signal.Ignore(syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT)
	_, env_exists := os.LookupEnv("HEROKU_REDIRECTED")
	if (!env_exists) && Exists(os.Getenv("LOCALAPPDATA")+"/heroku/client/bin/heroku.cmd") {
		os.Setenv("HEROKU_REDIRECTED", "1")
		startExe(os.Getenv("LOCALAPPDATA")+"/heroku/client/bin/heroku.exe", os.Args[1:])
	}
	_, env_exists = os.LookupEnv("HEROKU_BINPATH")
	if !env_exists {
		os.Setenv("HEROKU_BINPATH", dp0+"/heroku.exe")
	}
	if Exists(os.Getenv("HEROKU_BINPATH")) {
		startExe(dp0+"/../bin/node.exe", append([]string{dp0 + "/../bin/run"}, os.Args[1:]...))
	}
	cmd_content_bytes, _ := ioutil.ReadFile(dp0 + "/heroku.cmd")
	cmd_content := string(cmd_content_bytes)
	raw_path := strings.Split(cmd_content, "\"")[17]
	node_path := strings.Replace(raw_path, "\\", "/", -1)
	node_path = strings.Replace(node_path, "%LOCALAPPDATA%", os.Getenv("LOCALAPPDATA"), -1)
	if Exists(node_path) {
		startExe(node_path, append([]string{dp0 + "/../bin/run"}, os.Args[1:]...))
	}
	startExe("node", append([]string{dp0 + "/../bin/run"}, os.Args[1:]...))
}
