FROM golang:1.14
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go mod download
RUN go get -v github.com/cosmtrek/air
ENTRYPOINT air -c ./.air.toml
