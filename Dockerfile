FROM golang:1.13 

WORKDIR /go/src/app

COPY . .

RUN make

CMD ["/go/src/app/build/miner"]