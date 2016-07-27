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
http://www.devdungeon.com/content/working-files-go#create_empty_file
*/

/*
TODO: include GPS or any rough location
TODO: Send localhost name of device to server
*/

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"
	"../lib/"
)

func getArgs() (string, string) {
	var connType string = "tcp"            // or unix
	var connAddr string = "127.0.0.1:8088" // or /tmp/unix.sock

	flag.StringVar(&connType, "t", "tcp", "specify connection type.  defaults to tcp.")
	flag.StringVar(&connAddr, "a", "127.0.0.1:8088", "specify connection address.  defaults to 127.0.0.1:8088")
	flag.Parse()

	return connType, connAddr
}

func main() {
	var debug bool = true

	connType, connAddr := getArgs()

	c, err := net.Dial(connType, connAddr)
	if err != nil {
		panic(err)
	}
	defer c.Close()

	for {
		source := lib.GetSource()

		_, err := c.Write([]byte(lib.GetOutboundIP() + " Public: " + lib.Get_IP_WIMIA()))

		if debug == true {
			fmt.Println("Public: " + source)
			print(lib.GetOutboundIP() + "\n")
		}
		if err != nil {
			log.Fatal("write error:", err.Error())
			os.Exit(1)
		}

		time.Sleep(1800e9) // 1800 seconds
	}
}

