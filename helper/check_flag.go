package helper

import (
	"errors"
)

//Returns an error message when the input flag arugment is a default value.
func CheckFlags(function string, email string, password string) (err error) {
	switch {
	case function == "default_function":
		err = errors.New("Error: Function is missing!")
	case email == "default_email":
		err = errors.New("Error: Email is missing!")
	case password == "default_password":
		err = errors.New("Error: Password is missing!")
	default:
		err = nil
	}
	return
}