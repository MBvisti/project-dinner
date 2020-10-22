# Build stage
FROM golang:1.15-alpine AS build-env

ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

LABEL maintainer="Morten Vistisen vistisen@live.dk"

WORKDIR /app

ARG VERSION

COPY go.sum go.mod ./

RUN go mod download

COPY cmd ./cmd
COPY pkg ./pkg
COPY template ./template

RUN GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=${VERSION} -s -w" -a -o main cmd/server/main.go

FROM alpine
# to make the program have time zone data
RUN apk add --no-cache tzdata
COPY --from=build-env /app/main /
COPY --from=build-env /app/template/ /template/

CMD ["./main"]
