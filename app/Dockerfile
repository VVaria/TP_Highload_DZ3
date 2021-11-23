FROM golang:alpine

WORKDIR /go/src/app

ADD . .
RUN go build  -o main.go
EXPOSE 8080

CMD ["./main.go"]