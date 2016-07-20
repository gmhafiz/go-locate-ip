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