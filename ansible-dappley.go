package main 

import (
	"os"
	"fmt"
	"flag"
	"bufio"
	"strconv"
	"strings"
	"os/exec"
	"io/ioutil"
)

func main() {
	var function string
	flag.StringVar(&function, "function", "<Function Name>", "Name of the function that will be run.")
	flag.Parse()

	if function == "update" {
		update()
	} else if function == "terminate" {
		terminate()
	} else if function == "initialize" {
		initialize()
	} else {
		fmt.Println("Function Invalid!")
	}
}

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

	for i := 1; i <= 1; i++ {
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
				fmt.Println(instance_id)
				describe_instance := "aws ec2 describe-instances --instance-id " + instance_id
				args = strings.Split(describe_instance, " ")
				cmd := exec.Command(args[0], args[1:]...)
				output, err := cmd.CombinedOutput()
				if err != nil {
					fmt.Println(err)
				}
				fmt.Printf("%s\n", output)
				instance_ids += instance_id + "\n"
			}

			if strings.Contains(line, "PrivateIpAddress") {
				args := strings.Split(line, ": ")
				private_ip := strings.TrimLeft(strings.TrimRight(args[1], "\","), "\"")
				fmt.Println(private_ip)
				private_ips += private_ip + "\n"
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

func initialize() {
	fileName := "instance_ids"

	instance_byte, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println("Failed to read", fileName)
		return
	}

	scanner := bufio.NewScanner(strings.NewReader(string(instance_byte)))
	for i := 1; scanner.Scan() && i <= 1; i++ {
		instance_id := scanner.Text()
		initializing := true
		for initializing {
			status_file, err := os.Create("status")
			if err != nil {
				fmt.Println("Unable to create file!")
				return
			}

			terminate_instance := "aws ec2 describe-instance-status --instance-ids " + instance_id
			args := strings.Split(terminate_instance, " ")
			cmd := exec.Command(args[0], args[1:]...)
			output, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Println(err)
			}
			
			//fmt.Println(terminate_instance)
			//fmt.Printf("%s\n", output)

			_, err = status_file.WriteString(string(output))
			if err != nil {
				fmt.Println("Unable update status!")
				return
			}

			status_byte, err := ioutil.ReadFile("status")
			if err != nil {
				fmt.Println("Failed to check status")
				return
			}

			status_scanner := bufio.NewScanner(strings.NewReader(string(status_byte)))
			for status_scanner.Scan() {
				line := status_scanner.Text()

				if strings.Contains(line, "\"Status\":") {
					args := strings.Split(line, ": ")
					status := strings.TrimLeft(strings.TrimRight(args[1], "\""), "\"")
					fmt.Println(status)
					if status == "passed" {
						initializing = false
						break
					}
				}
			}
		}
	}
}

func terminate() {
	fileName := "instance_ids"
	
	instance_byte, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println("Failed to read", fileName)
		return
	}

	scanner := bufio.NewScanner(strings.NewReader(string(instance_byte)))
	for i := 1; scanner.Scan() && i <= 1; i++ {
		instance_id := scanner.Text()

		terminate_instance := "aws ec2 terminate-instances --instance-ids " + instance_id
		args := strings.Split(terminate_instance, " ")
		cmd := exec.Command(args[0], args[1:]...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("%s\n", output)
		fmt.Println(terminate_instance)
	}
}