package email

import (
	"github.com/heesooh/go-dappley-multi-server-tester/helper"
	"io/ioutil"
	"strings"
	"bufio"
)

//Create email with all failing cases
func ComposeEmail(fileNames []string) string {
	var file_description []string
	var emailContents string
	var emptySlice []string
	var task []string

	failingFiles := helper.IsFileFail(fileNames)
	if (len(failingFiles) == 0) {
		return "ALL TESTS PASS!"
	}
	files_results := make([][]string, len(failingFiles))

	for index := 0; index < len(failingFiles); index++ {
		curr_file := failingFiles[index]
		files_results[index] = append(files_results[index], curr_file)
		file_byte, err := ioutil.ReadFile(curr_file)
		if err != nil {
			files_results[index] = append(files_results[index], "Playbook: [" + strings.TrimLeft(strings.TrimRight(curr_file, ".txt"), "./test_results/") + "]")
			files_results[index] = append(files_results[index], "Failed to read playbook!")
			files_results[index] = append(files_results[index], "------------------------------------------------------")
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
				task = append(task, helper.Between(line, "[", "]"))
				continue
			}
			task = append(task, line)
			if line == "" {
				if helper.DoubleContains(task, "fatal:") {
					simplified_task := helper.Simplify(task)
					file_description = append(file_description, simplified_task...)
				}
				task = emptySlice
			}
		}
	}

	for _, result := range files_results {
		emailContents += "Playbook: [" + strings.TrimRight(result[0], ".txt")[15:] + "]\n"
		emailContents += "Failing Tasks: \n\n"
		for i, line := range result {
			if i == 0 {
				continue
			}
			emailContents += line + "\n"
		}
		emailContents += "------------------------------------------------------\n"
	}
	return emailContents
}