.PHONY: miner test image clean fmt

miner:
	go build -o build/miner

test:
	go test -cpu 1,4 -timeout 7m github.com/CommonProsperity/Miner/...

image:
	docker build -t commonprosperity/miner --build-arg GOPROXY=`go env GOPROXY` .

clean:
	rm -rf build/*

fmt:
	go fmt github.com/CommonProsperity/Miner