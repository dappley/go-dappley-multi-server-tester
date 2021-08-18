package helper

import (
	"strings"
)

//Creates a simpler email.
func simpleEmail(fileNames []string) string {
	var emailContents string
	failingFiles := IsFileFail(fileNames)
	if (len(failingFiles) == 0) {
		return "ALL TESTS PASS!"
	}
	emailContents += "Failing Playbook: \n"
	for _, fileName := range failingFiles {
		emailContents += "[" + strings.TrimRight(fileName, ".txt") + "] - Failing\n\n"
	}
	return emailContents
}

//Simplifies the data inside the email content.
func Simplify(task []string) []string {
	var fatal bool
	var simplified_task []string

	for i, _ := range task {
		if i == 0 || task[i] == "" {
			simplified_task = append(simplified_task, task[i])
			continue
		}
		if strings.Contains(task[i], "...ignoring") {
			continue
		}
		if strings.Contains(task[i], "skipping:") {
			continue
		}
		if strings.Contains(task[i], "fatal: ") {
			fatal = true
		}
		if strings.Contains(task[i], "ok: ") || strings.Contains(task[i], "changed: ") {
			fatal = false
		}
		if fatal {
			simplified_task = append(simplified_task, task[i])
		}
	}
	return simplified_task
}