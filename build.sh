#!/bin/bash

rm -f bin/*.exec

go build -o bin/server.exec cmd/server/main.go