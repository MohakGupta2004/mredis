package handlers

import (
	"io"
	"log"
	"net"
	"strings"

	"github.com/MohakGupta2004/mredis/internal/resp"
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
		parts := resp.ParseRequest(string(data))
		if len(parts) == 0 {
			log.Println("connection: command not found", string(data))
			return
		}
		command := strings.ToLower(parts[0])
		args := parts[1:]
		response := services.ExecuteCommand(command, args, store)
		conn.Write([]byte(response))
	}
}
