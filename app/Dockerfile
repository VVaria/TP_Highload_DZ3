FROM golang:alpine

WORKDIR /go/src/app

ADD . .
RUN go build  -o main.go
EXPOSE 80

CMD ["./main.go"]