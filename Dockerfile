FROM golang:1.14 AS builder

WORKDIR /go/src/app

COPY . .

RUN make

FROM ubuntu:20.04

WORKDIR /app

COPY --from=builder /go/src/app/build/miner .

CMD ["/app/miner"]