package internal

import (
	"log"
	"net"
	"os"
)

func PrintPenguin(conn net.Conn) {
	file, err := os.ReadFile("./welcome.txt")
	if err != nil {
		log.Fatalf("Could not open file: %v", err)
	}

	conn.Write(file)
}
