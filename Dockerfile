FROM golang:1.13.6-buster AS builder

WORKDIR /go/src/github.com/jdanper/logour
COPY . .

RUN go get -u github.com/golang/dep/...
RUN make build

FROM alpine:3.11.3

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /go/src/github.com/jdanper/logour/logour .

CMD ["./logour"]