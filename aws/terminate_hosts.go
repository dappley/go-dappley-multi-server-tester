package aws

import (
	"fmt"
	"log"
	"bufio"
	"strconv"
	"strings"
	"os/exec"
	"io/ioutil"
)

//Termiante all servers via aws cli command.
func Terminate_hosts(number string) {
	var updated_instance_list string
	var updated_host_list string
	fileName_1 := "instance_ids"
	fileName_2 := "hosts"

	to_terminate, err := strconv.Atoi(number)
	if err != nil {
		panic(err)
	}

	lines_to_remove := to_terminate * 2
	hosts_byte, err := ioutil.ReadFile(fileName_2)
	if err != nil {
		fmt.Println("Failed to read", fileName_2, "!")
		return
	}
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
	if err != nil {
		log.Fatalln(err)
	}

	instance_byte, err := ioutil.ReadFile(fileName_1)
	if err != nil {
		fmt.Println("Failed to read", fileName_1, "!")
		return
	}

	instance_scanner := bufio.NewScanner(strings.NewReader(string(instance_byte)))
	for instance_scanner.Scan() {
		instance_id := instance_scanner.Text()
		if to_terminate == 0 {
			updated_instance_list += instance_id + "\n"
			continue
		}
		terminate_instance := "aws ec2 terminate-instances --instance-ids " + instance_id
		args := strings.Split(terminate_instance, " ")
		cmd := exec.Command(args[0], args[1:]...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("%s\n", output)
		fmt.Println(terminate_instance)
		to_terminate -= 1
	}
	err = ioutil.WriteFile(fileName_1, []byte(updated_instance_list), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}