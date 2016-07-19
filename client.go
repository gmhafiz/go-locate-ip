package main

import (
	"log"
	"net"
	"time"
	"strings"
)

func main() {
	c, err := net.Dial("unix", "/tmp/echo.sock")
	if err != nil {
		panic(err)
	}
	defer c.Close()

	for {
		_, err := c.Write([]byte(GetOutboundIP()))
		if err != nil {
			log.Fatal("write error:", err)
			break
		}
		time.Sleep(10e9)  // 10 seconds
	}
}

// Get preferred outbound ip of this machine
func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().String()
	idx := strings.LastIndex(localAddr, ":")

	return localAddr[0:idx]
}