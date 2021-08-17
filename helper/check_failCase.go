package helper

import (
	"io/ioutil"
	"strings"
	"bufio"
)

//Check which given files contains failing test case
func IsFileFail(fileNames []string) []string {
	var failingFiles []string
	var curr_file string
	var scan_result bool

	for i := 0; i < len(fileNames); i++ {
		// Reset bool value
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
				if !(strings.Contains(line, "failed=0")) || !(strings.Contains(line, "ignored=0")) || !(strings.Contains(line, "unreachable=0")) {
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