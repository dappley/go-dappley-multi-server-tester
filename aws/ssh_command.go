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

//Prints out the ssh command for all servers.
func SSH_command(number string) {
	number_of_instances, err := strconv.Atoi(number)
	if err != nil { panic(err) }
	instance_byte, err := ioutil.ReadFile("instance_ids")
	if err != nil { log.Fatal("Failed to read instance_ids!") }
	scanner := bufio.NewScanner(strings.NewReader(string(instance_byte)))
	for i := 1; scanner.Scan() && i <= number_of_instances; i++ {
		describe_instance := "aws ec2 describe-instances --instance-ids " + scanner.Text()
		output := helper.ShellCommandExecuter(describe_instance)

		description_scanner := bufio.NewScanner(strings.NewReader(string(output)))
		for description_scanner.Scan() {
			line := description_scanner.Text()
			if strings.Contains(line, "\"PublicIpAddress\":") {
				public_ip := helper.TrimLeftRight(line, "\"", "\",")
				fmt.Println("ssh -i jenkins.pem ubuntu@" + public_ip)
				break
			}
		}
	}
}