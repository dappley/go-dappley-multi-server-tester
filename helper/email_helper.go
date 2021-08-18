package helper

import (
	"path/filepath"
	"net/mail"
	"strings"
	"os"
)

//Checks if slice contains the given value.
func contains(slice []string, val string) bool {
	for _, elem := range slice {
		if elem == val {
			return true
		}
	}
	return false
}

//Runs contains method twice.
func DoubleContains(slice []string, val string) bool {
	for _, elem := range slice {
		if strings.Contains(elem, val) {
			return true
		}
	}
	return false
}

//Get substring between two strings.
func Between(value string, a string, b string) string {
    posFirst := strings.Index(value, a)
    if posFirst == -1 {
        return ""
    }
    posLast := strings.Index(value, b)
    if posLast == -1 {
        return ""
    }
    posFirstAdjusted := posFirst + len(a)
    if posFirstAdjusted >= posLast {
        return ""
    }
    return value[posFirstAdjusted:posLast]
}

//Checks the validity of the email address.
func Valid_email(email string) bool {
    _, err := mail.ParseAddress(email)
    return err == nil
}

//Returns the slice containing name of all files in a drectory.
func AllFiles(directory string) []string {
    var files []string
    root := directory
    err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if path[len(path)-4:] == ".yml" || path[len(path)-4:] == ".txt" {
			files = append(files, "./" + path)
		}
        return nil
    })
    if err != nil {
        panic(err)
    }
	return files
}
