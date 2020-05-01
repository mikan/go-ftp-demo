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
		text = convertUnixLikeCommand(text)
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

func convertUnixLikeCommand(cmd string) string {
	if strings.HasPrefix(cmd, "cd") {
		return strings.Replace(cmd, "cd", "CWD", 1)
	}
	if strings.HasPrefix(cmd, "ls") {
		return strings.Replace(cmd, "ls", "NLST", 1)
	}
	if strings.HasPrefix(cmd, "dir") {
		return strings.Replace(cmd, "dir", "LIST", 1)
	}
	if strings.HasPrefix(cmd, "cat") {
		return strings.Replace(cmd, "cat", "RETR", 1)
	}
	if strings.HasPrefix(cmd, "rm") {
		return strings.Replace(cmd, "rm", "DELE", 1)
	}
	if strings.HasPrefix(cmd, "pwd") {
		return strings.Replace(cmd, "pwd", "PWD", 1)
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
