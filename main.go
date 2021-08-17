package main 

import (
	"os"
	"fmt"
	"flag"
	"path/filepath"
)

func main() {
	var number, function, senderEmail, senderPasswd string
	flag.StringVar(&number, "number", "999999", "Number of the ec2 instances to be terminated.")
	flag.StringVar(&function, "function", "<Function Name>", "Name of the function that will be run.")
	flag.StringVar(&senderEmail, "senderEmail", "<Sender Email>", "Email of the addressee.")
	flag.StringVar(&senderPasswd, "senderPasswd", "<Sender Password>", "Email password of the addressee.")
	flag.Parse()

	if function == "update" {
		update_hosts(number)

	} else if function == "initialize" {
		initialize_hosts(number)

	} else if function == "ssh_command" {
		ssh_command(number)

	} else if function == "update_address" {
		Update_address(allFiles("playbooks"))

	} else if function == "send_result" {
		SendTestResult(senderEmail, senderPasswd, allFiles("test_results"))
	
	} else if function == "terminate" {
		terminate_hosts(number)

	} else {
		fmt.Println("Function Invalid!")
	}
}

//----------Helper----------

func allFiles(directory string) []string {
    var files []string
    root := directory
    err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if path[len(path)-4:] == ".yml" || path[len(path)-4:] == ".txt" {
			files = append(files, "./" + path)
		}
        return nil
    })
    if err != nil {
        panic(err)
    }
	return files
}