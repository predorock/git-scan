package git

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"predorock/gitscan/logstreamer"
	"strings"
)

var logger = logstreamer.GetInstance()
var log = logger.Log

func checkCommand(command string) bool {
	return strings.HasPrefix(command, "git ")
}

func checkCommandWithMsg(command string, msg string) bool {
	if !checkCommand(command) {
		log.Fatalln(msg)
		return false
	}
	return true
}

func checkCommandOrExit(command string) {
	if !checkCommandWithMsg(command, fmt.Sprintf("Not a git command: <%s>\n", command)) {
		os.Exit(1)
	}
}

func ExecGitCommand(command string, dir string) {
	checkCommandOrExit(command)
	var args = strings.Split(command, " ")
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = logstreamer.NewLogstreamer(log, "stdout", path.Base(dir), false)
	cmd.Stderr = logstreamer.NewLogstreamer(log, "stderr", path.Base(dir), false)
	cmd.Dir = dir
	cmd.Run()
}

func UpdateCommand(folder string) {
	log.Printf("Updating repo %s\n", path.Base(folder))
	ExecGitCommand("git pull -f", folder)
}

func CreateGitCommand(command string, desc func(string) string) func(string) {
	checkCommandOrExit(command)
	return func(folder string) {
		log.Println(desc(folder))
		ExecGitCommand(command, folder)
	}

}
