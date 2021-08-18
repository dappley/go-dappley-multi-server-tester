package aws

import (
	"github.com/heesooh/go-dappley-multi-server-tester/helper"
	"io/ioutil"
	"strings"
	"strconv"
	"errors"
	"bufio"
	"fmt"
	"log"
)

//Runs until all servers are initialized.
func Initialize_hosts(number string) {
	//number is type string because of Jenkins pipeline
	instances_to_initialize, err := strconv.Atoi(number)
	if err != nil { panic(err) }
	fileName := "instance_ids"
	instance_byte, err := ioutil.ReadFile(fileName)
	if err != nil { log.Fatal("Failed to read", fileName) }

	scanner := bufio.NewScanner(strings.NewReader(string(instance_byte)))
	for i := 1; scanner.Scan() && i <= instances_to_initialize; i++ {
		instance_id := scanner.Text()
		initializing := true
		fmt.Println("Initializing " + instance_id + "...")
		for initializing {
			initialize_instance := "aws ec2 describe-instance-status --instance-ids " + instance_id
			output := helper.ShellCommandExecuter(initialize_instance)
			status_scanner := bufio.NewScanner(strings.NewReader(string(output)))
			for status_scanner.Scan() {
				line := status_scanner.Text()
				if strings.Contains(line, "\"InstanceStatuses\":") {
					status := helper.TrimLeftRight(line, "\"", "\"")
					if status == "[]" {
						err := errors.New("Instance " + instance_id + "has been termianted!")
						panic(err)
					}
				}
				if strings.Contains(line, "\"Status\":") {
					status := helper.TrimLeftRight(line, "\"", "\"")
					if status == "passed" {
						initializing = false
						fmt.Println("Instance " + instance_id + " initialized!")
						break
					}
				}
			}
		}
	}
}