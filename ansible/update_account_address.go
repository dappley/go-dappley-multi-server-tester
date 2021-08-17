package ansible 

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
	var peer_ids []string

	for i := 1; i <= 5; i++ {
		account_address, err := ioutil.ReadFile("./accounts/node" + strconv.Itoa(i) + "/account_address.txt")
		if err != nil {
			fmt.Println("Failed to read account" + strconv.Itoa(i) + "'s address!")
			continue
		}
		account_addresses = append(account_addresses, string(account_address))
	}

	for i := 1; i <= 5; i++ {
		peer_id, err := ioutil.ReadFile("./accounts/node" + strconv.Itoa(i) + "/node" + strconv.Itoa(i) + "_peerID.txt")
		if err != nil {
			fmt.Println("Failed to read node" + strconv.Itoa(i) + "'s peerID!")
			continue
		}
		peer_ids = append(peer_ids, string(peer_id))
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
			} else if strings.Contains(line, "peer_ID_1") {
				updated_line := strings.ReplaceAll(line, "peer_ID_1", peer_ids[0]) // + "\n"
				updated_playbook += updated_line
			} else if strings.Contains(line, "peer_ID_2") {
				updated_line := strings.ReplaceAll(line, "peer_ID_2", peer_ids[1]) // + "\n"
				updated_playbook += updated_line
			} else if strings.Contains(line, "peer_ID_3") {
				updated_line := strings.ReplaceAll(line, "peer_ID_3", peer_ids[2]) // + "\n"
				updated_playbook += updated_line
			} else if strings.Contains(line, "peer_ID_4") {
				updated_line := strings.ReplaceAll(line, "peer_ID_4", peer_ids[3]) // + "\n"
				updated_playbook += updated_line
			} else if strings.Contains(line, "peer_ID_5") {
				updated_line := strings.ReplaceAll(line, "peer_ID_5", peer_ids[4]) // + "\n"
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