package aws

import (
	"github.com/heesooh/go-dappley-multi-server-tester/helper"
	"io/ioutil"
	"strconv"
	"strings"
	"bufio"
	"fmt"
	"log"
)

//Termiante all servers via aws cli command.
func Terminate_hosts(number string) {
	var updated_instance_list string
	var updated_host_list string
	fileName_1 := "instance_ids"
	fileName_2 := "hosts"

	to_terminate, err := strconv.Atoi(number)
	if err != nil { panic(err) }

	lines_to_remove := to_terminate * 2
	hosts_byte, err := ioutil.ReadFile(fileName_2)
	if err != nil { log.Fatal("Failed to read", fileName_2, "!") }

	host_scanner := bufio.NewScanner(strings.NewReader(string(hosts_byte)))
	for host_scanner.Scan() {
		line := host_scanner.Text()
		if lines_to_remove == 0 {
			updated_host_list += line + "\n"
			continue
		}
		lines_to_remove -= 1
	}
	err = ioutil.WriteFile(fileName_2, []byte(updated_host_list), 0644)
	if err != nil { log.Fatalln(err) }

	instance_byte, err := ioutil.ReadFile(fileName_1)
	if err != nil { log.Fatal("Failed to read", fileName_1, "!") }

	instance_scanner := bufio.NewScanner(strings.NewReader(string(instance_byte)))
	for instance_scanner.Scan() {
		instance_id := instance_scanner.Text()
		if to_terminate == 0 {
			updated_instance_list += instance_id + "\n"
			continue
		}
		terminate_instance := "aws ec2 terminate-instances --instance-ids " + instance_id
		output := helper.ShellCommandExecuter(terminate_instance)
		fmt.Printf("%s\n", output)
		fmt.Println(terminate_instance)
		to_terminate -= 1
	}
	err = ioutil.WriteFile(fileName_1, []byte(updated_instance_list), 0644)
	if err != nil { log.Fatalln(err) }
}