# Locate My Laptop

## Introduction

Find the IP address of your (lost) devices by making it send both public and local IP address.
Motivation

## Usage

Run the server first before running the client

```go
go run lml_server.go
```

```go
go run lml_client.go
```

## Configuration

The line `CONN_ADDR = "127.0.0.1:8088"` has to be changed to whatever your server IP address is.

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
- systemd service that starts upon network connection