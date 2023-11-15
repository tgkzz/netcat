package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"net-cat/internal"
	asd "net-cat/internal"
)

func main() {
	port := ""
	if len(os.Args) == 1 {
		port = asd.DefaultPort
	} else if len(os.Args) == 2 {
		port = os.Args[1]
	} else {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}
	listen, err := net.Listen("tcp", "localhost:"+port)
	if err != nil {
		fmt.Println("Error when starting the server:", err)
		return
	}
	fmt.Printf("Listening on the port :%s\n", port)
	defer listen.Close()

	ch := make(chan asd.Message)
	go asd.Broadcast(ch)

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Print("qwe")
			log.Fatal(err)
		}
		asd.Mutex.Lock()
		//asd.Connections++
		if asd.Connections >= asd.MaxConnections {
			asd.Mutex.Unlock()
			conn.Write([]byte("server is busy. please try again later"))
			conn.Close()
			continue
		}
		asd.Connections++
		asd.Mutex.Unlock()
		fmt.Println(internal.Clients)
		//go asd.HandleConnection(conn, ch)
		go func() {
			asd.HandleConnection(conn, ch)
			asd.Mutex.Lock()
			asd.Connections-- // Decrement the connection count when the client disconnects
			asd.Mutex.Unlock()
		}()
	}
}
