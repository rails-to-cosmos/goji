FROM golang

WORKDIR /app

ADD main.go main.go
ADD go.mod go.mod
ADD go.sum go.sum

RUN go build

ENTRYPOINT ["/app/goji"]
