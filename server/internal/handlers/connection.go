package handlers

import (
	"io"
	"log"
	"net"

	"github.com/MohakGupta2004/mredis/internal/services"
	"github.com/MohakGupta2004/mredis/internal/utils"
)

func HandleConn(conn net.Conn, store *services.StoreType) {
	defer conn.Close()
	for {
		data, err := utils.ParseReq(conn)
		if err != nil {
			if err == io.EOF {
				return
			}
			log.Println("connection: couldn't able to parse the req", err.Error())
			return
		}
		opcode, args := services.ParseOpcode(string(data))
		command := services.ParseCommand(opcode, args)
		if command == "" {
			log.Println("connection: command not found", string(data))
			return
		}
		response := services.ExecuteCommand(command, args, store)
		conn.Write([]byte(response))
	}
}
