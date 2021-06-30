package main 

import (
	"log"
	"fmt"
	"bufio"
	"strings"
	"strconv"
	"io/ioutil"
)

func Update_address(playbooks []string) {
	var account_addresses []string
	for i := 1; i <= 5; i++ {
		account_address, err := ioutil.ReadFile("./accounts/node" + strconv.Itoa(i) + "/account_address.txt")
		if err != nil {
			fmt.Println("Failed to read account" + strconv.Itoa(i) + "'s address!")
			continue
		}
		account_addresses = append(account_addresses, string(account_address))
	}

	for _, playbook := range playbooks {
		var updated_playbook string
		playbook_byte, err := ioutil.ReadFile(playbook)
		if err != nil {
			fmt.Println("Failed to read " + playbook)
			return
		}

		scanner := bufio.NewScanner(strings.NewReader(string(playbook_byte)))
		for scanner.Scan() {
			line := scanner.Text()

			if strings.Contains(line, "account_1_address") {
				updated_line := strings.ReplaceAll(line, "account_1_address", account_addresses[0][:34]) + "\n"
				updated_playbook += updated_line

			} else if strings.Contains(line, "account_2_address") {
				updated_line := strings.ReplaceAll(line, "account_2_address", account_addresses[1][:34]) + "\n"
				updated_playbook += updated_line

			} else if strings.Contains(line, "account_3_address") {
				updated_line := strings.ReplaceAll(line, "account_3_address", account_addresses[2][:34]) + "\n"
				updated_playbook += updated_line

			} else if strings.Contains(line, "account_4_address") {
				updated_line := strings.ReplaceAll(line, "account_4_address", account_addresses[3][:34]) + "\n"
				updated_playbook += updated_line

			} else if strings.Contains(line, "account_5_address") {
				updated_line := strings.ReplaceAll(line, "account_5_address", account_addresses[4][:34]) + "\n"
				updated_playbook += updated_line

			} else {
				updated_playbook += line + "\n"

			}
		}
		err = ioutil.WriteFile(playbook, []byte(updated_playbook), 0644)
		if err != nil {
			log.Fatalln(err)
		}
	}
}