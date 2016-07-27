# Locate My Laptop

## Introduction

Find the IP address of your (lost) devices by making it send both public and local IP address.
Motivation

## Installation

```shell
go get github.com/gmhafiz/locateMyLaptop
```

## Usage

Run the server first before running the client. 

```go
go run lml_server.go
```

```go
go run lml_client.go
```

## Configuration

Flags are available for `lml-client.go`. 

If you are running the server on the internet, please check if the port `8088` is not firewalled and yu have he correct IP address.

```
-t <type of connection> // defaults to tcp
-a <address:port>       // defaults to 127.0.0.1:8808
```

## Autostart

It is better to compile it though. `go install lml_[server|client].go` Since there are two mains in both .go files.

Using systemd:

```
[Unit]g
Description=Send both local and public IP address to an external server

[Service]
ExecStart=/usr/bin/lml_client
Restart=on-abort

[Install]
WantedBy=multi-user.target

```

## TODO

- Send rough GPS location
- Send localhost name of device to server
- systemd service that starts upon network connection

## Bug

- Closing (Ctrl + C) client will close server as well
- Check if log file exists or not