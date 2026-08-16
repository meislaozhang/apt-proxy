.PHONY: test vet build clean

test:
	go test ./...

vet:
	go vet ./...

build:
	go build -o bin/apt-server ./cmd/apt-server
	go build -o bin/apt-client ./cmd/apt-client

clean:
	rm -rf bin
