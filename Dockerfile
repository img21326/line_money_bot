FROM golang:1.15

WORKDIR /go/src/app
COPY . .

RUN go mod download
RUN go build .
CMD ["./moneybot"]