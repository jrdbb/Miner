FROM golang:1.14 

WORKDIR /go/src/app

COPY . .

RUN make

CMD ["/go/src/app/build/miner"]