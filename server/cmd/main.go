package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/MohakGupta2004/mredis/internal/handlers"
	"github.com/MohakGupta2004/mredis/internal/services"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	store := services.NewStore()
	defer stop()
	lc := net.ListenConfig{}
	l, err := lc.Listen(ctx, "tcp", ":6379")
	if err != nil {
		log.Fatalln("error: couldn't able to connect to port 6379: ", err.Error())
	}
	log.Println("server started")
	defer l.Close()

	go func() {
		<-ctx.Done()
		l.Close()
	}()

	for {
		conn, err := l.Accept()
		if err != nil {
			if ctx.Err() != nil {
				log.Println("server shutting down")
				return
			}
			log.Fatalln("error: couldn't able to connect to port 6379: ", err.Error())
			continue
		}
		log.Println("client connected")
		go handlers.HandleConn(conn, store)
	}
}
