package helper

import (
	"os/exec"
	"strings"
	"fmt"
)

//Trims both left and right and returns the trimmed result
func TrimLeftRight(line string, left string, right string) (info string) {
	args := strings.Split(line, ": ")
	info = strings.TrimLeft(strings.TrimRight(args[1], right), left)
	return
}

//Executes the shell command and returns the output
func ShellCommandExecuter(command string) (output []byte) {
	args := strings.Split(command, " ")
	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil { fmt.Println(err) }
	return
}