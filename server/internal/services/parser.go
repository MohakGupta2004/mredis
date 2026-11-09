package services

import (
	"log"
	"strings"
)

func ParseOpcode(data string) (byte, []string) {
	args := strings.Split(data, "\r\n")
	operation := args[0][0]
	return operation, args
}

func ParseCommand(opcode byte, args []string) string {
	switch opcode {
	case '$':
		data := args[1]
		return data
	case '*':
		if len(args) < 3 {
			log.Println("parseCommand: invalid command format")
			return ""
		}
		command := args[2]
		command = strings.ToLower(command)
		return command
	}

	return ""
}
