package aws

import (
	"os"
	"fmt"
	"bufio"
	"strconv"
	"strings"
	"io/ioutil"
)

//Adds the server information to the hosts and instance_ids file.
func Update_hosts(number string) {
	instances_to_update, err := strconv.Atoi(number)
	if err != nil {
		panic(err)
	}
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

	for i := 1; i <= instances_to_update; i++ {
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