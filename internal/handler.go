package internal

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

func HandleConnection(conn net.Conn, ch chan<- Message) {
	defer conn.Close()
	PrintPenguin(conn)
	name := PrintName(conn)


	if _, ok := Clients[name]; ok {
		conn.Write([]byte("this nickname already use\n"))
		conn.Close()
		return
	} else {
		Clients[name] = conn
	}


	for _, message := range MessageHistory {
		conn.Write([]byte(message + "\n"))
	}

	onetime := time.Now().Format("2006-01-02 15:04:05")
	connMessage := Message{
		info: "\n" + name + " has joined our chat ... \n",
		body: "[" + onetime + "]",
		from: conn,
	}

	ch <- connMessage

	for {
		timeNow := time.Now().Format("2006-01-02 15:04:05")
		terminal := "[" + timeNow + "]" + "[" + name + "]" + ":"
		conn.Write([]byte(terminal))
		msg, _, err := bufio.NewReader(conn).ReadLine()
		if IsKeys(msg) {
			conn.Write([]byte("Unavailable character, you will be kicked"))
			conn.Close()
		}
		if err != nil {
			Mutex.Lock()
			Connections--
			delete(Clients, name)
			Mutex.Unlock()
			fmt.Printf("%s disconnected\n", name)
			connMessage := Message{
				info: "\n" + name + " has left our chat...\n",
				body: "[" + timeNow + "]",
				from: conn,
			}
			ch <- connMessage
			break
		}
		if IsPrintable(string(msg)) {
			connMessage := Message{
				body: "\n" + terminal + string(msg) + "\n" + "[" + timeNow + "]",
				from: conn,
			}
			ch <- connMessage
			MessageHistory = append(MessageHistory, terminal+string(msg))
		}
	}
}


func PrintName(conn net.Conn) string {
	for {
		conn.Write([]byte("[ENTER YOUR NAME]:"))
		takeName, _, err := bufio.NewReader(conn).ReadLine()
		if err != nil {
			log.Fatal(err)
		}

		if IsKeys(takeName) {
			conn.Write([]byte("Unavailable character, you will be kicked"))
			conn.Close()
			return ""
		}

		name := string(takeName)

		if len(name) < 2 || len(name) > 32 {
			conn.Write([]byte("Name must be between 2 and 32 characters\n"))
			continue
		}

		if IsUsed(name) {
			conn.Write([]byte("Nickname is already in use\n"))
			continue
		}

		if !IsPrintable(name) {
			conn.Write([]byte("Nickname can consist only of Latin letters. You will be kicked from the chat\n"))
			continue
		}

		fmt.Println(name + " connected")
		return name
	}
}


func IsPrintable(msg string) bool {
	printableFlag := false
	for _, char := range msg {
		if char != ' ' && char != '\t' && char != '\n' && char != '\r' {
			printableFlag = true
			if char < 32 || char > 126 {
				return false
			}
		}
	}
	return printableFlag
}

func IsKeys(bytes []byte) bool {
	for i := 0; i < len(bytes); i++ {
		if bytes[i] == 27 && bytes[i+1] == 91 && i+1 < len(bytes) {
			return true
		}
	}
	return false
}

func IsUsed(str string) bool {
	for name := range Clients {
		if name == str {
			return true
		}
	}

	return false
}
