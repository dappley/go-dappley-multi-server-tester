package main

import (
	"fmt"
	"bufio"
	"strconv"
	"strings"
	"os/exec"
	"io/ioutil"
)

//Prints out the ssh command for all servers
func ssh_command(number string) {
	number_of_instances, err := strconv.Atoi(number)
	if err != nil {
		panic(err)
	}
	instance_byte, err := ioutil.ReadFile("instance_ids")
	if err != nil {
		fmt.Println("Failed to read instance_ids!")
		return
	}
	scanner := bufio.NewScanner(strings.NewReader(string(instance_byte)))
	for i := 1; scanner.Scan() && i <= number_of_instances; i++ {
		instance_id := scanner.Text()
		describe_instance := "aws ec2 describe-instances --instance-ids " + instance_id
		args := strings.Split(describe_instance, " ")
		cmd := exec.Command(args[0], args[1:]...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println(err)
		}
		description_scanner := bufio.NewScanner(strings.NewReader(string(output)))
		for description_scanner.Scan() {
			line := description_scanner.Text()
			if strings.Contains(line, "\"PublicIpAddress\":") {
				public_ip_args := strings.Split(line, ": ")
				public_ip := strings.TrimLeft(strings.TrimRight(public_ip_args[1], "\","), "\"")
				fmt.Println("ssh -i jenkins.pem ubuntu@" + public_ip)
				break
			}
		}
	}
}