package main

import (
	"fmt"
	brokerManager "github.com/JosephCottingham/mqtt_interface_cli/brokerManager"
	// "time"
	"golang.org/x/term"
	"os"
	"syscall"
)

func main() {
	fmt.Println("start")
	passwordCollect := true
	for passwordCollect {
		password := PasswordPrompt("What is your password?")
		brokerData, err := brokerManager.ReadBrokerData(password)
		if err != nil {
			fmt.Println(err)
		} else {
			passwordCollect = false
			StartShell(brokerData, password)
		}
	}
}

func PasswordPrompt(label string) string {
	var s string
	for {
		fmt.Fprint(os.Stderr, label+" ")
		b, _ := term.ReadPassword(int(syscall.Stdin))
		s = string(b)
		if s != "" {
			break
		}
	}
	fmt.Println()
	return s
}
