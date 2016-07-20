package main

/*
Credits:
https://www.socketloop.com/tutorials/golang-convert-http-response-body-to-string
https://www.socketloop.com/tutorials/golang-find-ip-address-from-string
https://www.socketloop.com/tutorials/golang-convert-http-response-body-to-string
http://stackoverflow.com/questions/106179/regular-expression-to-match-dns-hostname-or-ip-address
*/

/*
TODO: include GPS or any rough location
*/

import (
	"log"
	"net"
	"time"
	"strings"
	"net/http"
	"os"
	"fmt"
	"io/ioutil"
	"regexp"
)

const (
	CONN_TYPE = "tcp"		// unix or tcp
	CONN_ADDR = "127.0.0.1:8088"	// /tmp/echo.sock or addr:port
)

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

func get_ip_wimia() string {
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


func get_source() string {
	//response, err := http.Get("https://www.google.com.au/search?q=whattismyipaddress&oq=whattismyipaddress&aqs=chrome..69i57j0l5.7070j0j7&sourceid=chrome&ie=UTF-8")
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

func find_ip(data string) string {

	re, err := regexp.Compile(`[0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5].`)
	if err != nil {
		fmt.Println("ip not found")
		os.Exit(1)
	}
	ip := re.Find([]byte(data))

	return string(ip)
}

//func GetPublicIP() string {
//	response, err := http.Get("http://www.whatismyipaddress.net")
//	if err != nil {
//		os.Stderr.WriteString(err.Error())
//		os.Stderr.WriteString("\n")
//		os.Exit(1)
//	}
//	defer response.Body.Close()
//	html_data, err := ioutil.ReadAll(response.Body)
//	if err != nil {
//		fmt.Println(err)
//		os.Exit(1)
//	}

	//ip, err := regexp.Match("ip", html_data)
	//if err != nil {
	//	fmt.Printf("IP not found?!\n")
	//	return
	//}
	//var ip string
	//ip = findIP(string(html_data))

	//fmt.Printf(ip)
	//return ip
	//return "ip"
//}

func findIP(input string) string {
	//numBlock := "(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])"
	//numBlock := "([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])"
	//numBlock := "([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5]).([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5]).([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5]).([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])"
	//class := "_h4c"
	//numBlock := "(?:[0-9]{1,3})[0-9]{1,3}"
	//regexPattern := numBlock + "\\." + numBlock + "\\." + numBlock + "\\." + numBlock

	//find_class := regexp.MustCompile(class)
	//regEx := regexp.MustCompile(class)

	//r, _ := regexp.Compile(`14.202.167.25`)
	//res := r.FindAllString(input, -1)
	//fmt.Printf("%v", res)

	//return regEx.FindString(res)
	//return res
	return "res"
}

func main() {
	var debug bool = true
	c, err := net.Dial(CONN_TYPE, CONN_ADDR)
	if err != nil {
		panic(err)
	}
	defer c.Close()

	for {
		source := get_source()

		_, err := c.Write([]byte(GetOutboundIP() + " " + get_ip_wimia()))

		if debug == true {
			fmt.Println("Public IP: " + source)
			print(GetOutboundIP() + "\n")
		}
		if err != nil {
			log.Fatal("write error:", err.Error())
			os.Exit(1)
		}
		time.Sleep(1800e9)  // 1800 seconds
	}
}
