.PHONY: all
all: build

.PHONY: build
build:
	go build -o zhero main.go

# Cross-compile for Raspberry Pi Zero W (armv6l)
.PHONY: build-rpi-zero
build-rpi-zero:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -o build/zheroapp-rpi-zero main.go

# --- RPI Deployment ---
SERVER ?= $(SSH_SERVER)
PORT ?= $(SSH_PORT)
add-rpi-key:
	ssh-keygen -t rsa -b 1024 -f ~/.ssh/rpi-zero -N ""
	ssh-copy-id $(SERVER)

copy-rpi:
	ssh -p $(PORT) $(SERVER) "mkdir -p /tmp/zhero"
	scp -P $(PORT) -r build/zheroapp-rpi-zero config.yaml template $(SERVER):/tmp/zhero

# --- Android Packaging (GoMobile) ---

# Defin- GOBIN for consistent gomobile command access
GOBIN := $(shell go env GOBIN)
ifeq ($(GOBIN),)
    GOBIN := $(shell go env GOPATH)/bin
endif

# Initialize gomobile environment
.PHONY: gomobile-init
gomobile-init:
	go install golang.org/x/mobile/cmd/gomobile@latest
	gomobile init

.PHONY: build-android-lib
build-android-lib:
	# CGO_ENABLED=1 is required for gomobile.
	CGO_ENABLED=1 gomobile bind -tags "android fts5" -target=android -androidapi 21 -o zheroapp.aar ./server
	mv zheroapp.aar zhero-android-app/app/libs/zheroapp.aar
	mv zheroapp-sources.jar zhero-android-app/app/libs/zheroapp-sources.jar


# TODO copy config + rdf_schema + zhero.db + template/ files to Android

# --- Android Application Building and Deployment ---
# ANDROID_APP_DIR = zhero-android-app
# ADB ?= adb # Ensure adb is in your PATH, or set this variable to its full path (e.g., /path/to/android/sdk/platform-tools/adb)
#
# # Build the Android APK
# .PHONY: build-android-apk
# build-android-apk: build-android-lib
# 	@echo "Building Android APK..."
# 	@cd $(ANDROID_APP_DIR) && ./gradlew assembleDebug
#
# # Push the Android APK to a connected device
# .PHONY: push-android-apk
# push-android-apk: build-android-apk
# 	@echo "Installing Android APK on device..."
# 	@$(ADB) install -r $(ANDROID_APP_DIR)/app/build/outputs/apk/debug/app-debug.apk

