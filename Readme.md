# go-locate-ip

## Introduction

Find both internal and external IP address as well as rough GPS location of your (lost) devices by making it send data to your server at a regular interval.

#### Motivation
[![ DEF CON 18 - Zoz - Pwned By The Owner: What Happens When You Steal A Hacker's Computer ](http://img.youtube.com/vi/Jwpg-AwJ0Jc/0.jpg)](https://www.youtube.com/watch?v=Jwpg-AwJ0Jc)

## Installation

```shell
go get github.com/gmhafiz/go-locate-ip
```

## Usage

1. Build program
```go
go build main.go
````

Run the server first before running the client. 

```go
./go-locate-ip
```

2. And then run client

```go
./go-locate-ip -m client
```

## Configuration

Flags are available for this program. 

If you are running the server on the internet, please check if the port `8088` is not firewalled and you have the correct IP address.

```
-m <type of connection>     // defaults to server. Options are server and client 
-a <address>                // defaults to 0.0.0.0
-p <port>                   // default to 8088
-t <intervals in minutes>   // defaults to 10 minutes
```

Example:
Connect laptop to a server at 216.58.199.78 port 8000 and updates every 30 minutes
```bash
./main -m client -a 216.58.199.78 -p 8000 -t 30
```

## Autostart

It is better to compile it though. `go install main.go`. Make sure you have `$GOBIN` set up and is included in your `$PATH`

Using systemd:

```
[Unit]
Description=Send both local and public IP address and GPS location to an external server

[Service]
ExecStart=/usr/bin/go-locate-ip
Restart=on-abort

[Install]
WantedBy=multi-user.target

```

## TODO

- ~~Send rough GPS location~~
- Send localhost name of device to server
- systemd service that starts upon network connection

## Bug

- Closing (Ctrl + C) client will close server as well
- ~~Check if log file exists or not~~