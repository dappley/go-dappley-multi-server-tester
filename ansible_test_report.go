package main

import (
	"fmt"
	"bufio"
	"strings"
	"io/ioutil"
	"gopkg.in/gomail.v2"
)

//Send email
func SendTestResult(recipient string, senderEmail string, senderPasswd string) {
	fileNames := []string{"setup.txt", "single_transaction_no_tip.txt"}
	emailContents := detailedEmail(fileNames)

	fmt.Println(emailContents)

	send(recipient, emailContents, fileNames, senderEmail, senderPasswd)
}

func send(recipient string, emailBody string, attachment []string, senderEmail string, senderPasswd string) {
	//send the email
	mail := gomail.NewMessage()
	mail.SetHeader("From", senderEmail)
	mail.SetHeader("To",   recipient)
	mail.SetHeader("Subject", "Ansible Test Result")
	mail.SetBody("text", emailBody)
	for _, fileName := range attachment {
		mail.Attach(fileName)
	}

	deliver := gomail.NewDialer("smtp.gmail.com", 587, senderEmail, senderPasswd)

	if err := deliver.DialAndSend(mail); err != nil {
		fmt.Println("Failed to send email!")
		panic(err)
	}
}

//Create simpler email
func simpleEmail(fileNames []string) string {
	failingFiles := isFileFail(fileNames)
	if (len(failingFiles) == 0) {
		return "ALL TESTS PASS!"
	}

	var emailContents string
	emailContents += "Failing Playbook: \n"
	for _, fileName := range failingFiles {
		emailContents += "[" + strings.TrimRight(fileName, ".txt") + "] - Failing\n\n"
	}

	return emailContents
}

func simplify(task []string) []string {
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

	fmt.Println(task)

	fmt.Println(simplified_task)

	return simplified_task
}

//Create email with all failing cases
func detailedEmail(fileNames []string) string {
	failingFiles  := isFileFail(fileNames)
	files_results := make([][]string, len(failingFiles))

	if (len(failingFiles) == 0) {
		return "ALL TESTS PASS!"
	}

	var file_description []string
	var task []string
	var emailContents string
	var emptySlice []string

	for index := 0; index < len(failingFiles); index++ {
		curr_file := failingFiles[index]
		files_results[index] = append(files_results[index], curr_file)

		file_byte, err := ioutil.ReadFile(curr_file)
		if err != nil {
			files_results[index] = append(files_results[index], "Read File")
			files_results[index] = append(files_results[index], "Unable to read" + curr_file + "!")
			continue
		}

		scanner := bufio.NewScanner(strings.NewReader(string(file_byte)))
		for i := 0; scanner.Scan(); i++ {
			if i == 0 || i == 1 || i == 2 {
				continue
			}

			line := scanner.Text()

			if strings.Contains(line, "PLAY RECAP") {
				files_results[index] = append(files_results[index], file_description...)
				file_description = emptySlice
				break
			}

			if strings.Contains(line, "TASK [") {
				task = append(task, between(line, "[", "]"))
				continue
			}

			task = append(task, line)

			if line == "" {
				if doubleContains(task, "fatal:") {
					simplified_task := simplify(task)
					file_description = append(file_description, simplified_task...)
				}
				task = emptySlice
			}
		}
	}

	for _, result := range files_results {
		emailContents += "Playbook: [" + strings.TrimRight(result[0], ".txt") + "]\n"
		emailContents += "Failing Tasks: \n\n"
		for i, line := range result {
			if i == 0 {
				continue
			}
			emailContents += line + "\n"
		}
		emailContents += "---------------------------\n"
	}

	return emailContents
}

//Check which given files contains failing test case
func isFileFail(fileNames []string) []string {
	var failingFiles []string
	var curr_file string
	var scan_result bool

	for i := 0; i < len(fileNames); i++ {
		if curr_file != fileNames[i] {
			curr_file = fileNames[i]
			scan_result = false
		}

		file_byte, err := ioutil.ReadFile(curr_file)
		if err != nil {
			failingFiles = append(failingFiles, curr_file)
			continue
		}
		
		scanner := bufio.NewScanner(strings.NewReader(string(file_byte)))
		for scanner.Scan() {
			line := scanner.Text()
			if line == "" {
				continue
			}

			if scan_result {
				if !(strings.Contains(line, "failed=0")) || !(strings.Contains(line, "ignored=0")) {
					if contains(failingFiles, curr_file) {
						continue
					} else {
						failingFiles = append(failingFiles, curr_file)
						continue
					}
				}
			}
			if strings.Contains(line, "PLAY RECAP") {
				scan_result = true
			}
		}
	}
	return failingFiles
}

//Checks if slice contains the given value
func contains(slice []string, val string) bool {
	for _, elem := range slice {
		if elem == val {
			return true
		}
	}
	return false
}

func doubleContains(slice []string, val string) bool {
	for _, elem := range slice {
		if strings.Contains(elem, val) {
			return true
		}
	}
	return false
}

//Get substring between two strings.
func between(value string, a string, b string) string {
    posFirst := strings.Index(value, a)
    if posFirst == -1 {
        return ""
    }
    posLast := strings.Index(value, b)
    if posLast == -1 {
        return ""
    }
    posFirstAdjusted := posFirst + len(a)
    if posFirstAdjusted >= posLast {
        return ""
    }
    return value[posFirstAdjusted:posLast]
}