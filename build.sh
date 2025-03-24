#!/bin/bash

rm -f bin/*

go build -o bin/server.exec cmd/server/main.go