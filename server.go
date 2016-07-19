package main

import (
	"log"
	"net"
	"time"
)

func what_is_the_ip(c net.Conn) {
	for {
		buf := make([]byte, 512)
		nr, err := c.Read(buf)
		if err != nil {
			return
		}

		ip := buf[0:nr]
		time := time.Now().Format(time.RFC850)
		data := time + " " + string(ip)
		println("Server got:", string(data))
		//_, err = c.Write(data)
		//if err != nil {
		//	log.Fatal("Write: ", err)
		//}
	}
}

func main() {
	l, err := net.Listen("unix", "/tmp/echo.sock")
	if err != nil {
		log.Fatal("listen error:", err)
	}

	for {
		incoming, err := l.Accept()
		if err != nil {
			log.Fatal("accept error:", err)
		}

		go what_is_the_ip(incoming)
	}
}