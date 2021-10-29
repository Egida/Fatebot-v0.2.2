#!/bin/bash
read -p 'Please enter your payload name: ' pName
go build -ldflags "-s -w" -o $pName main.go
echo Build success...