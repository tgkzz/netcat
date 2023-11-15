package internal

func Broadcast(ch <-chan Message) {
	for {
		msg := <-ch
		Mutex.Lock()
		for name, conn := range Clients {
			if conn == msg.from {
				continue
			}
			if msg.info == "" {
				conn.Write([]byte(msg.body))
				bname := "[" + name + "]" + ":"
				conn.Write([]byte(bname))
			} else {
				conn.Write([]byte(msg.info))
				aname := msg.body + "[" + name + "]" + ":"
				conn.Write([]byte(aname))
			}
		}
		Mutex.Unlock()
	}
}
