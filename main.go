package main 

import (
	"github.com/heesooh/go-dappley-multi-server-tester/ansible"
	"github.com/heesooh/go-dappley-multi-server-tester/helper"
	"github.com/heesooh/go-dappley-multi-server-tester/email"
	"github.com/heesooh/go-dappley-multi-server-tester/aws"
	"flag"
	"fmt"
	"log"
)

func main() {
	var number, function, senderEmail, senderPasswd string
	flag.StringVar(&number, "number", "999999", "Number of the ec2 instances to be terminated.")
	flag.StringVar(&function, "function", "default_function", "Name of the function that will be run.")
	flag.StringVar(&senderEmail, "senderEmail", "default_email", "Email of the addressee.")
	flag.StringVar(&senderPasswd, "senderPasswd", "default_password", "Email password of the addressee.")
	flag.Parse()

	err := helper.CheckFlags(function, senderEmail, senderPasswd)
	if err != nil {
		log.Fatal(err)
		return
	}

	if function == "update" {
		aws.Update_hosts(number)
	} else if function == "initialize" {
		aws.Initialize_hosts(number)
	} else if function == "ssh_command" {
		aws.SSH_command(number)
	} else if function == "update_address" {
		ansible.Update_address(helper.AllFiles("playbooks"))
	} else if function == "send_result" {
		email.SendTestResult(senderEmail, senderPasswd, helper.AllFiles("test_results"))
	} else if function == "terminate" {
		aws.Terminate_hosts(number)
	} else {
		fmt.Println("Function Invalid!")
	}
}