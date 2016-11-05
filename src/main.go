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
https://golang.org/pkg/flag/
*/

import (
	"fmt"
	"log"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"time"
	"flag"
	"strings"
)

const (
	CONN_TYPE = "tcp"   // unix or tcp
	CONN_PORT = ":8088" // any port >= 1024
)

var (
	newFile		*os.File
	err		error
	filename	string = "log.txt"
	mode		string
	serverAddr	string
	serverPort	int = 8088
	pollTime	int = 60
)

func runClient() {
	fmt.Println("Running client")
	fmt.Println("Connecting to " + serverAddr)

	c, err := net.Dial("tcp", serverAddr)
	if err != nil {
		panic(err)
	}
	defer c.Close()
	for {
		var internalIP string = GetInternalIP()
		var externalIP string = GetExternalIP()

		_, err := c.Write([]byte("\n\tExternal IP: " + externalIP + "\n\tInternal ip: " + internalIP))
		if err != nil {
			log.Fatal("write error:", err)
			break
		}
		time.Sleep(10e9)  // 10 seconds
	}
}

// Get preferred outbound ip of this machine
func GetInternalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().String()
	idx := strings.LastIndex(localAddr, ":")

	return "Local: " + localAddr[0:idx]
}

// Internal IP
func GetExternalIP() string {
	response, err := http.Get("http://ipv4bot.whatismyipaddress.com")
	if err != nil {
		log.Fatal("404 not found", err.Error())
		os.Exit(1)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal("http read error", err.Error())
		os.Exit(1)
	}
	src := string(body)

	return string(src)
}

func runServer() {
	listen, err := net.Listen(CONN_TYPE, CONN_PORT)
	if err != nil {
		log.Fatal("listen error:", err.Error())
	}
	defer listen.Close()

	currTime := time.Now().Format(time.RFC850)
	fmt.Println(currTime + ": Listening to incoming connections...")

	for {
		incoming, err := listen.Accept()
		if err != nil {
			log.Fatal("accept error:", err.Error())
		}

		go getIPFromClient(incoming)
	}
}

func getIPFromClient(conn net.Conn) {
	for {
		buf := make([]byte, 512)
		length, err := conn.Read(buf)
		if err != nil {
			log.Fatal("Error connecting: " + err.Error())
			os.Exit(1)
		}

		ip := buf[0:length]
		currTime := time.Now().Format(time.RFC850)
		data := currTime + " " + string(ip)
		println("Server got:", string(data))
		appendToLog(data)  // Write incoming message into a log file.
	}
	conn.Close()
}

func appendToLog(src string) {
	createFile(filename)
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err.Error())
	}
	defer file.Close()

	if _, err = file.WriteString(src + "\n"); err != nil {
		panic(err.Error())
	}
}

func createFile(filename string) {
	newFile, err = os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}

	newFile.Close()
}



func init() {
	flag.StringVar(&mode, "m", "server", "Mode: -m server|client")
	flag.StringVar(&serverAddr, "a", "0.0.0.0:8088", "Server address: -a 127.0.0.1")
	flag.IntVar(&serverPort, "p", 8088, "Server port: -p 8088")
	flag.IntVar(&pollTime, "t", 60, "Poll interval (seconds): -t 60")
}

func main() {

	flag.Parse()
	fmt.Println(mode)

	if mode == "server" {
		runServer()
		return
	} else if mode == "client" {
		runClient()
		return
	} else {
		fmt.Println("No recognisable flag")
	}
}
