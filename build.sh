#!/bin/bash

packr2

GOOS=windows GOARCH=amd64 go build -ldflags "-w -s" -o result/webssh_windows_amd64.exe .
GOOS=windows GOARCH=386 go build -ldflags "-w -s" -o result/webssh_windows_386.exe .
GOOS=linux GOARCH=amd64 go build -ldflags "-w -s" -o result/webssh_linux_amd64 .
GOOS=linux GOARCH=arm64 go build -ldflags "-w -s" -o result/webssh_linux_arm64 .
GOOS=darwin GOARCH=amd64 go build -ldflags "-w -s" -o result/webssh_darwin_amd64 .

packr2 clean
