# Build stage
FROM golang:1.14-alpine AS build-env

ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

LABEL maintainer="Morten Vistisen vistisen@live.dk"

WORKDIR /app

ARG VERSION

COPY go.sum .
COPY go.mod .

RUN go mod download
RUN go get -u github.com/cosmtrek/air

COPY cmd ./cmd
COPY pkg ./pkg
COPY .air.toml .
COPY template ./template

# TODO: should probably look into using go install here
# RUN GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=${VERSION} -s -w" -a -o main cmd/server/main.go

FROM alpine
# to make the program have time zone data
RUN apk add --no-cache tzdata
COPY --from=build-env /app /
RUN go get -u github.com/cosmtrek/air
# COPY --from=build-env /app/template/daily_recipe_email.html /template/
RUN ls
ENTRYPOINT air -c ./.air.toml
