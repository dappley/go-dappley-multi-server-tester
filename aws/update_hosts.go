package aws

import (
	"github.com/heesooh/go-dappley-multi-server-tester/helper"
	"io/ioutil"
	"strconv"
	"strings"
	"bufio"
	"log"
	"os"
)

//Adds the server information to the hosts and instance_ids file.
func Update_hosts(number string) {
	instances_to_update, err := strconv.Atoi(number)
	if err != nil { panic(err) }

	//Create txt files for server info
	host_file, err := os.Create("hosts")
	if err != nil { log.Fatal("Unable to create file!") }
	id_file, err := os.Create("instance_ids")
	if err != nil { log.Fatal("Unable to create file!") }

	for i := 1; i <= instances_to_update; i++ {
		var private_ips, instance_ids string
		fileName := "node" + strconv.Itoa(i) + ".txt"
		node_byte, err := ioutil.ReadFile(fileName)
		if err != nil { log.Fatal("Failed to read", fileName) }

		scanner := bufio.NewScanner(strings.NewReader(string(node_byte)))
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "InstanceId") {
				instance_id := helper.TrimLeftRight(line, "\"", "\",")
				instance_ids += instance_id + "\n"
			}
			if strings.Contains(line, "PrivateIpAddress") {
				private_ip := helper.TrimLeftRight(line, "\"", "\",")
				private_ips += "[NODE" + strconv.Itoa(i) + "]\n" + private_ip + "\n"
				break
			}
		}

		_, err = id_file.WriteString(instance_ids)
		if err != nil { log.Fatal("Unable to write on file!") }
		_, err = host_file.WriteString(private_ips)
		if err != nil { log.Fatal("Unable to write on file!") }
	}
}