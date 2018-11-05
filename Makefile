VERSION  := $$(git describe --tags --always)
TARGET   := marathon_exporter
TEST     ?= ./...

default: test build

test:
	go test -v -run=$(RUN) $(TEST)

build: clean
	go build -v -o bin/$(TARGET)

release: clean
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build \
		-a -tags netgo \
		-ldflags "-X main.Version=$(VERSION)" \
		-o bin/$(TARGET) .
	docker build -t crerwin/$(TARGET):$(VERSION) .

clean:
	rm -rf bin/
