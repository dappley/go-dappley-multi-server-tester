package main 

import (
	"os"
	"fmt"
	"log"
	"flag"
	"bufio"
	"errors"
	"strconv"
	"strings"
	"os/exec"
	"io/ioutil"
	"path/filepath"
)

func main() {
	var number, function, recipient, senderEmail, senderPasswd string
	flag.StringVar(&number, "number", "999999", "Number of the ec2 instances to be terminated.")
	flag.StringVar(&function, "function", "<Function Name>", "Name of the function that will be run.")
	flag.StringVar(&recipient, "recipient", "<Recipient Email>", "Email of the recipient.")
	flag.StringVar(&senderEmail, "senderEmail", "<Sender Email>", "Email of the addressee.")
	flag.StringVar(&senderPasswd, "senderPasswd", "<Sender Password>", "Email password of the addressee.")
	flag.Parse()

	if function == "update" {
		update()

	} else if function == "initialize" {
		initialize()

	} else if function == "ssh_command" {
		ssh_command()

	} else if function == "update_address" {
		Update_address(allFiles("playbooks"))

	} else if function == "send_result" {
		SendTestResult(recipient, senderEmail, senderPasswd, allFiles("testresults"))
	
	} else if function == "terminate" {
		terminate(number)

	} else {
		fmt.Println("Function Invalid!")
	}
}

//Adds the server information to the hosts and instance_ids file
func update() {
	//Create txt files for server info
	host_file, err := os.Create("hosts")
	if err != nil {
		fmt.Println("Unable to create file!")
		return
	}

	id_file, err := os.Create("instance_ids")
	if err != nil {
		fmt.Println("Unable to create file!")
		return
	}

	for i := 1; i <= 5; i++ {
		var private_ips, instance_ids string
		fileName := "node" + strconv.Itoa(i) + ".txt"
		
		node_byte, err := ioutil.ReadFile(fileName)
		if err != nil {
			fmt.Println("Failed to read", fileName)
			return
		}

		scanner := bufio.NewScanner(strings.NewReader(string(node_byte)))
		for scanner.Scan() {
			line := scanner.Text()

			if strings.Contains(line, "InstanceId") {
				args := strings.Split(line, ": ")
				instance_id := strings.TrimLeft(strings.TrimRight(args[1], "\","), "\"")
				instance_ids += instance_id + "\n"
			}

			if strings.Contains(line, "PrivateIpAddress") {
				args := strings.Split(line, ": ")
				private_ip := strings.TrimLeft(strings.TrimRight(args[1], "\","), "\"")
				private_ips += "[NODE" + strconv.Itoa(i) + "]\n" + private_ip + "\n"
				break
			}
		}

		_, err = host_file.WriteString(private_ips)
		if err != nil {
			fmt.Println("Unable to write on file!")
			return
		}

		_, err = id_file.WriteString(instance_ids)
		if err != nil {
			fmt.Println("Unable to write on file!")
			return
		}
	}
}

//Runs until all servers are initialized
func initialize() {
	fileName := "instance_ids"
	instance_byte, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println("Failed to read", fileName)
		return
	}

	scanner := bufio.NewScanner(strings.NewReader(string(instance_byte)))
	for i := 1; scanner.Scan() && i <= 5; i++ {
		instance_id := scanner.Text()
		initializing := true
		fmt.Println("Initializing " + instance_id + "...")
		for initializing {
			initialize_instance := "aws ec2 describe-instance-status --instance-ids " + instance_id
			args := strings.Split(initialize_instance, " ")
			cmd := exec.Command(args[0], args[1:]...)
			output, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Println(err)
			}
			
			status_scanner := bufio.NewScanner(strings.NewReader(string(output)))
			for status_scanner.Scan() {
				line := status_scanner.Text()

				if strings.Contains(line, "\"InstanceStatuses\":") {
					args := strings.Split(line, ": ")
					status := strings.TrimLeft(strings.TrimRight(args[1], "\""), "\"")
					if status == "[]" {
						err := errors.New("Instance " + instance_id + "has been termianted!")
						panic(err)
					}
				}

				if strings.Contains(line, "\"Status\":") {
					args := strings.Split(line, ": ")
					status := strings.TrimLeft(strings.TrimRight(args[1], "\""), "\"")
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

//Termiante all servers via aws cli command
func terminate(number string) {
	var updated_instance_list string
	fileName := "instance_ids"

	to_terminate, err := strconv.Atoi(number)
	if err != nil {
		panic(err)
	}

	instance_byte, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println("Failed to read", fileName, "!")
		return
	}

	scanner := bufio.NewScanner(strings.NewReader(string(instance_byte)))
	for scanner.Scan() {
		instance_id := scanner.Text()
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

	err = ioutil.WriteFile(fileName, []byte(updated_instance_list), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}

//Prints out the ssh command for all servers
func ssh_command() {	
	instance_byte, err := ioutil.ReadFile("instance_ids")
	if err != nil {
		fmt.Println("Failed to read instance_ids!")
		return
	}

	scanner := bufio.NewScanner(strings.NewReader(string(instance_byte)))
	for i := 1; scanner.Scan() && i <= 5; i++ {
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
    // for _, file := range files {
    //     fmt.Println(file)
	// }
	
	return files
}