package main

/*
Credits / References:
https://gist.github.com/jniltinho/9788121
https://stackoverflow.com/questions/5885486/how-do-i-get-the-current-timestamp-in-go
https://systembash.com/a-simple-go-tcp-server-and-tcp-client/
https://www.socketloop.com/tutorials/golang-convert-http-response-body-to-string
https://www.socketloop.com/tutorials/golang-find-ip-address-from-string
https://www.socketloop.com/tutorials/golang-convert-http-response-body-to-string
http://stackoverflow.com/questions/106179/regular-expression-to-match-dns-hostname-or-ip-address
http://www.devdungeon.com/content/working-files-go#create_empty_file
http://www.devdungeon.com/content/ip-geolocation-go
https://golang.org/pkg/flag/
https://stackoverflow.com/questions/12518876/how-to-check-if-a-file-exists-in-go
https://stackoverflow.com/questions/31786215/can-command-line-flags-in-go-be-set-to-mandatory
*/

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
	"flag"
	"strings"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"strconv"
)

const (
	CONN_TYPE = "tcp"
	CONN_PORT = ":8088" // any port >= 1024
)

var (
	newFile		*os.File
	filename	string = "log.txt"
	mode		string
	serverAddr	string
	serverPort	int = 8088
	pollMinutes	int = 10  // 10 minutes

	err      error
	geo      GeoIP
	response *http.Response
	body     []byte
)

type GeoIP struct {
	// The right side is the name of the JSON variable
	Ip          string  `json:"ip"`
	CountryCode string  `json:"country_code"`
	CountryName string  `json:"country_name"`
	RegionCode  string  `json:"region_code"`
	RegionName  string  `json:"region_name"`
	City        string  `json:"city"`
	Zipcode     string  `json:"zipcode"`
	Lat         float32 `json:"latitude"`
	Lon         float32 `json:"longitude"`
	MetroCode   int     `json:"metro_code"`
	AreaCode    int     `json:"area_code"`
}

/*
Geolocation
 */
func GeoLookup(address string) GeoIP {
	// Use freegeoip.net to get a JSON response
	// There is also /xml/ and /csv/ formats available
	response, err = http.Get("https://freegeoip.net/json/" + address)
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()

	// response.Body() is a reader type. We have
	// to use ioutil.ReadAll() to read the data
	// in to a byte slice(string)
	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}

	// Unmarshal the JSON byte slice to a GeoIP struct
	err = json.Unmarshal(body, &geo)
	if err != nil {
		fmt.Println(err)
	}

	return geo
}

/*
======
Client
======
 */
func runClient() {
	fmt.Println("Running on client mode.")
	fmt.Println("Connecting to " + serverAddr)

	c, err := net.Dial("tcp", serverAddr)
	if err != nil {
		panic(err)
	}
	defer c.Close()
	for {
		var internalIP string = GetInternalIP()
		var externalIP string = GetExternalIP()

		//fmt.Println(GeoLookup(externalIP).CountryName)

		geo := GeoLookup(externalIP)

		_, err := c.Write([]byte("\n\tExternal IP: " + externalIP +
			"\n\tInternal ip: " + internalIP +
			"\n\tCountry: " + geo.CountryName +
			"\n\tRegion: " + geo.RegionName +
			"\n\tCity: " + geo.City +
			"\n\tZipcode: " + geo.Zipcode +
			"\n\tLatitude: " + strconv.FormatFloat(float64(geo.Lat), 'f', -1, 64) +
			"\n\tLongitude: " + strconv.FormatFloat(float64(geo.Lon), 'f', -1, 64)))
		if err != nil {
			log.Fatal("write error:", err)
			break
		}
		//time.Sleep(pollMinutes * time.Minute)
		time.Sleep(time.Duration(pollMinutes) * time.Minute)
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

	return localAddr[0:idx]
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

/*
======
Server
======
 */
func runServer() {
	fmt.Println("Running on server mode")

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
			log.Fatal("Error connecting" + err.Error())
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

/*
=========
File IO
=========
 */
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
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		newFile, err = os.Create(filename)
		if err != nil {
			log.Fatal(err)
		}
		newFile.Close()
	}
}

/*
===============
Initialisation
===============
 */
func init() {
	flag.StringVar(&mode, "m", "server", "Mode: -m server|client")
	flag.StringVar(&serverAddr, "a", "0.0.0.0:8088", "Server address: -a 127.0.0.1")
	flag.IntVar(&serverPort, "p", serverPort, "Server port: -p 8088")
	flag.IntVar(&pollMinutes, "t", pollMinutes, "Poll interval (minutes): -t 10")
}

/*
============
Entry
============
 */
func main() {
	flag.Parse()
	if mode == "server" || mode == "s" {
		runServer()
		return
	} else if mode == "client" || mode == "c" {
		runClient()
		return
	} else {
		fmt.Println("No recognisable flag")
	}
}
