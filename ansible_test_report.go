package main

import (
	"fmt"
	"bufio"
	"strings"
	//"strconv"
	"io/ioutil"
)

//Send email
func SendTestResult() {
	fileNames := []string{"setup.txt", "single_transaction_no_tip.txt", "Unreadable.txt"}
	emailContents := composeEmail(fileNames)

	fmt.Println(emailContents)
}

//Create email
func composeEmail(fileNames []string) string {
	failingFiles  := isFileFail(fileNames)
	files_results := make([][]string, len(failingFiles))
	var file_description []string
	var task []string
	var emailContents string
	var emptySlice []string

	for index := 0; index < len(failingFiles); index++ {
		curr_file := failingFiles[index]
		file_byte, err := ioutil.ReadFile(curr_file)
		if err != nil {
			files_results[index] = append(files_results[index], curr_file)
			files_results[index] = append(files_results[index], "Unable to read file!")
			continue
		}

		scanner := bufio.NewScanner(strings.NewReader(string(file_byte)))
		for i := 0; scanner.Scan(); i++ {
			if i == 0 || i == 1 || i == 2 {
				continue
			}

			line := scanner.Text()

			if strings.Contains(line, "PLAY RECAP") {
				files_results[index] = file_description
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
					file_description = append(file_description, task...)
				}
				task = emptySlice
			}
		}
	}

	fmt.Println(files_results[0], "\n\n")
	fmt.Println(files_results[1], "\n\n")
	fmt.Println(files_results[2], "\n\n")

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

			if scan_result {
				if !(strings.Contains(line, "failed=0")) {
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