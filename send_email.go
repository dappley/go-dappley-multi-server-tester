package main

import (
	"fmt"
	"log"
	"bufio"
	"strings"
	"net/mail"
	"io/ioutil"
	"gopkg.in/gomail.v2"
)

func SendTestResult(senderEmail string, senderPasswd string, test_results []string) {
	emailContents := composeEmail(test_results)
	send(emailContents, senderEmail, senderPasswd)
}

func send(emailBody string, senderEmail string, senderPasswd string) {
	var recipients []string

	file_byte, err := ioutil.ReadFile("recipients.txt")
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(strings.NewReader(string(file_byte)))
	for scanner.Scan() {
		line := scanner.Text()
		if !valid_email(line) {
			fmt.Println("Invalid email address: \"" + line + "\"")
			continue
		}
		recipients = append(recipients, line)
	}
	//send the email
	mail := gomail.NewMessage()
	mail.SetHeader("From", senderEmail)
	addresses := make([]string, len(recipients))
	for i, recipient := range recipients {
		addresses[i] = mail.FormatAddress(recipient, "")
	}
	mail.SetHeader("To", addresses...)
	mail.SetHeader("Subject", "Ansible Test Result")
	mail.SetBody("text", emailBody)
	mail.Attach("test_results.zip")

	deliver := gomail.NewDialer("smtp.gmail.com", 587, senderEmail, senderPasswd)

	if err := deliver.DialAndSend(mail); err != nil {
		fmt.Println("Failed to send email!")
		panic(err)
	}
}

//----------Helper----------

//Checks the validity of the email address
func valid_email(email string) bool {
    _, err := mail.ParseAddress(email)
    return err == nil
}