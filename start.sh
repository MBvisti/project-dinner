#!/bin/sh

export $1
export $2
export $3
export $4
export $5
export $6
export $7

kill -9 `lsof -nP -i4TCP:5000 | grep LISTEN | tr -s ' ' | cut -d' ' -f2`
go vet ./...
go fmt ./...
go run cmd/server/main.go
