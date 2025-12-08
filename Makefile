.PHONY: all
all: build

.PHONY: build
build:
	go build -o zhero main.go

# Cross-compile for Raspberry Pi Zero W (armv6l)
.PHONY: build-rpi-zero
build-rpi-zero:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -o zhero-rpi-zero main.go

SERVER ?= $(SSH_SERVER)
PORT ?= $(SSH_PORT)
add-rpi-key:
	ssh-keygen -t rsa -b 1024 -f ~/.ssh/rpi-zero -N ""
	ssh-copy-id $(SERVER)

SERVER ?= $(SSH_SERVER)
PORT ?= $(SSH_PORT)
copy-rpi:
	ssh -p $(PORT) $(SERVER) "mkdir -p /tmp/zhero"
	scp -P $(PORT) -r zhero-rpi-zero config.yaml template $(SERVER):/tmp/zhero
