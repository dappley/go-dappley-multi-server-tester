package main 

import (
	"os"
	"fmt"
	"bufio"
	"strconv"
	"strings"
	"os/exec"
	"io/ioutil"
)

func main() {
	//Create txt file for server info
	file, err := os.Create("nodes.txt")
	if err != nil {
		fmt.Println("Unable to create file!")
		return
	}

	//Add the server info to the new txt file in the below format
	// Name:               <INSTANCE NAME>
	// Private ip address: <PRIVATE  IP>
	// Instance id:        <INSTANCE ID>
	for i := 1; i <= 5; i++ {
		var nodeName, privateIP, instanceID string
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
				instanceID = "Instance id:        " + instance_id + "\n"
				describe_instance := "aws ec2 describe-instances --instance-id " + instance_id
				args = strings.Split(describe_instance, " ")
				cmd := exec.Command(args[0], args[1:]...)
				output, err := cmd.CombinedOutput()
				if err != nil {
					fmt.Println(err)
				}
				fmt.Printf("%s\n", output)
			}

			if strings.Contains(line, "PrivateIpAddress") {
				args := strings.Split(line, ": ")
				private_ip := strings.TrimRight(args[1], ",")
				fmt.Println(private_ip)
				privateIP = "Private ip address: " + private_ip + "\n\n"
				break
			}
		}

		nodeName = "Name:               " + fileName + "\n"
		info := nodeName + instanceID + privateIP
		_, err = file.WriteString(info)
		if err != nil {
			fmt.Println("Unable to write on file!")
			return
		}
	}
}

