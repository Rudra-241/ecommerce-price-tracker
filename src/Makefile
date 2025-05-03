SRC := cmd/main.go
all: build
build:
	@mkdir -p bin
	@go build -o bin/app $(SRC)
run: build
	bin/app
clean:
	@rm -rf bin

.PHONY: all build run clean
