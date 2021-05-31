package main 

import (
	"os"
	"log"
	"fmt"
	// "bytes"
	"bufio"
	"strings"
	"strconv"
	"io/ioutil"
)

func Update_address() {
	//playbooks := []string{"accounts_generator.yml", "data.yml", "multi_transaction_no_tip.yml", "multi_transaction_with_tip.yml", "send.yml", "sendFromMiner.yml", "setup.yml", "single_transaction_no_tip.yml", "single_transaction_with_tip.yml", "smart_contract_gas_1.yml", "smart_contract_gas_2.yml", "update_seed_port.yml"}
	playbooks := []string{"multi_transaction_with_tip.yml"}
	var account_addresses []string

	for i := 1; i <= 5; i++ {
		account_address, err := ioutil.ReadFile("../go-dappley-ansible-accounts/node" + strconv.Itoa(i) + "/account_address.txt")
		if err != nil {
			fmt.Println("Failed to read account" + strconv.Itoa(i) + "'s address!")
			continue
		}
		account_addresses = append(account_addresses, string(account_address))
	}

	for _, playbook := range playbooks {
		// var playbook_updated []byte

		playbook_byte, err := ioutil.ReadFile(playbook)
		if err != nil {
			fmt.Println("Failed to read " + playbook)
			return
		}

		// playbook_updated = bytes.Replace(playbook_byte, []byte("account_1_address"), []byte(account_addresses[0]), -1)
		// playbook_updated = bytes.Replace(playbook_byte, []byte("account_2_address"), []byte(account_addresses[1]), -1)
		// playbook_updated = bytes.Replace(playbook_byte, []byte("account_3_address"), []byte(account_addresses[2]), -1)
		// playbook_updated = bytes.Replace(playbook_byte, []byte("account_4_address"), []byte(account_addresses[3]), -1)
		// playbook_updated = bytes.Replace(playbook_byte, []byte("account_5_address"), []byte(account_addresses[4]), -1)
		

		// if err = ioutil.WriteFile("test.txt", playbook_updated, 0666); err != nil {
		// 	fmt.Println(err)
		// 	return
		// }

		file, err := os.Create("test.yml")

		if err != nil {
			log.Fatal(err)
		}

		defer file.Close()

		fmt.Println("done")

		scanner := bufio.NewScanner(strings.NewReader(string(playbook_byte)))
		for scanner.Scan() {
			line := scanner.Text()

			// fmt.Println(line)

			// fmt.Println(strings.ReplaceAll(line, "account_1_address", account_addresses[0]))
			// fmt.Println(strings.ReplaceAll(line, "account_2_address", account_addresses[1]))
			// fmt.Println(strings.ReplaceAll(line, "account_3_address", account_addresses[2]))
			// fmt.Println(strings.ReplaceAll(line, "account_4_address", account_addresses[3]))
			// fmt.Println(strings.ReplaceAll(line, "account_5_address", account_addresses[4]))

			if strings.Contains(line, "account_1_address") {
				//fmt.Println(strings.ReplaceAll(line, "account_1_address", account_addresses[0]))
				_, err2 := file.WriteString(strings.ReplaceAll(line, "account_1_address", account_addresses[0]))
				if err2 != nil {
					log.Fatal(err2)
				}
			} else if strings.Contains(line, "account_2_address") {
				//fmt.Println(strings.ReplaceAll(line, "account_2_address", account_addresses[1]))
				_, err2 := file.WriteString(strings.ReplaceAll(line, "account_2_address", account_addresses[1]))
				if err2 != nil {
					log.Fatal(err2)
				}
			} else if strings.Contains(line, "account_3_address") {
				// fmt.Println(strings.ReplaceAll(line, "account_3_address", account_addresses[2]))
				_, err2 := file.WriteString(strings.ReplaceAll(line, "account_3_address", account_addresses[2]))
				if err2 != nil {
					log.Fatal(err2)
				}
			} else if strings.Contains(line, "account_4_address") {
				// fmt.Println(strings.ReplaceAll(line, "account_4_address", account_addresses[3]))
				_, err2 := file.WriteString(strings.ReplaceAll(line, "account_4_address", account_addresses[3]))
				if err2 != nil {
					log.Fatal(err2)
				}
			} else if strings.Contains(line, "account_5_address") {
				// fmt.Println(strings.ReplaceAll(line, "account_5_address", account_addresses[4]))
				_, err2 := file.WriteString(strings.ReplaceAll(line, "account_5_address", account_addresses[4]))
				if err2 != nil {
					log.Fatal(err2)
				}
			} else {
				_, err2 := file.WriteString(line)
				if err2 != nil {
					log.Fatal(err2)
				}
			}
		}
	}
}