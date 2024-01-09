package pkg

import "fmt"

func ErrorHandling(err error) {
	fmt.Println("[ERROR]: ", err.Error())
}

func StringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}
