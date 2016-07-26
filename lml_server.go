package main

/*
Credits:
https://gist.github.com/jniltinho/9788121
https://stackoverflow.com/questions/5885486/how-do-i-get-the-current-timestamp-in-go
https://systembash.com/a-simple-go-tcp-server-and-tcp-client/
https://www.socketloop.com/tutorials/golang-convert-http-response-body-to-string
https://www.socketloop.com/tutorials/golang-find-ip-address-from-string
https://www.socketloop.com/tutorials/golang-convert-http-response-body-to-string
http://stackoverflow.com/questions/106179/regular-expression-to-match-dns-hostname-or-ip-address
*/

import (
	"net"
	"time"
	"log"
	"os"
	"fmt"
)

const (
	CONN_TYPE = "tcp"	// unix or tcp
	CONN_PORT = ":8088"	// any port >= 1024
)

func appendToLog(src string) {
	file, err := os.OpenFile("log.txt", os.O_APPEND | os.O_WRONLY, 0600)
	if err != nil {
		panic(err.Error())
	}
	defer file.Close()
	if _, err = file.WriteString(src + "\n"); err != nil {
		panic(err.Error())
	}
}

func what_is_the_ip(conn net.Conn) {
	for {
		buf := make([]byte, 512)
		len, err := conn.Read(buf)
		if err != nil {
			log.Fatal("Error connecting: " + err.Error())
			os.Exit(1)
		}

		ip := buf[0:len]
		time := time.Now().Format(time.RFC850)
		data := time + " " + string(ip)
		println("Server got:", string(data))
		// Write incoming message into a log file.
		appendToLog(data)
	}
	conn.Close()
}

func main() {
	listen, err := net.Listen(CONN_TYPE, CONN_PORT)
	if err != nil {
		log.Fatal("listen error:", err.Error())
	}
	defer listen.Close()

	time := time.Now().Format(time.RFC850)
	fmt.Println(time + ": Listening to incoming connections...")

	for {
		incoming, err := listen.Accept()
		if err != nil {
			log.Fatal("accept error:", err.Error())
		}

		go what_is_the_ip(incoming)
	}
}