package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var dataCommands = []string{"RETR", "NLST", "LIST", "STOR", "APPE"}
var commandMap = map[string]string{"cd": "CMD", "ls": "NLST", "dir": "LIST", "cat": "RETR", "rm": "DELE", "pwd": "PWD"}

func main() {
	host := flag.String("h", "localhost", "hostname")
	port := flag.Int("port", 21, "port number")
	user := flag.String("u", "anonymous", "username")
	pass := flag.String("p", "anonymous@example.com", "password")
	debug := flag.Bool("d", false, "enable debug log")
	flag.Parse()
	if len(*host) == 0 {
		flag.Usage()
		os.Exit(2)
	}
	client, err := NewFTPClient(*host, *port)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer client.Close()
	if *debug {
		client.SetLogger(log.New(os.Stdout, "[FTP] ", 0))
	}
	if err := client.Login(*user, *pass); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	input := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for input.Scan() {
		text := input.Text()
		if len(text) == 0 {
			continue
		}
		if text == "exit" || text == "quit" {
			break
		}
		text = convertUserFriendlyCommand(text)
		if isDataCmd(text) {
			msg, err := client.DataCmd(text)
			if err != nil {
				os.Exit(1)
			}
			fmt.Println(msg)
		} else {
			code, msg, err := client.Cmd(text)
			if err != nil {
				os.Exit(1)
			}
			fmt.Printf("FTP %d\n%s\n", code, msg)
		}
		fmt.Print("> ")
	}
}

func convertUserFriendlyCommand(cmd string) string {
	for k, v := range commandMap {
		if strings.HasPrefix(cmd, k) {
			return strings.Replace(cmd, k, v, 1)
		}
	}
	return cmd
}

func isDataCmd(text string) bool {
	for _, cmd := range dataCommands {
		if strings.HasPrefix(strings.ToUpper(text), cmd) {
			return true
		}
	}
	return false
}
