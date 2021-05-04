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
	file, err := os.Create("hosts")
	if err != nil {
		fmt.Println("Unable to create file!")
		return
	}

	for i := 1; i <= 5; i++ {
		var private_ips string
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
			}

			if strings.Contains(line, "PrivateIpAddress") {
				args := strings.Split(line, ": ")
				private_ip := strings.TrimLeft(strings.TrimRight(args[1], "\","), "\"")
				fmt.Println(private_ip)
				private_ips += private_ip + "\n"
				break
			}
		}

		_, err = file.WriteString(private_ips)
		if err != nil {
			fmt.Println("Unable to write on file!")
			return
		}
	}
}

