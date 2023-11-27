# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=ActionPadServer
BINARY_UNIX=$(BINARY_NAME)_unix

all: build
build:
		$(GOBUILD) -o $(BINARY_NAME) -v
test:
		$(GOTEST) -v ./...
clean:
		$(GOCLEAN)
		rm -f $(BINARY_NAME)
		rm -f $(BINARY_UNIX)
run:
		$(GOBUILD) -o $(BINARY_NAME) -v ./...
		./$(BINARY_NAME)
deps:
		$(GOGET) github.com/gorilla/mux
		$(GOGET) github.com/go-vgo/robotgo
		$(GOGET) github.com/skratchdot/open-golang/open
		$(GOGET) github.com/sqweek/dialog
		$(GOGET) github.com/ActionPad/systray
		$(GOGET) github.com/skip2/go-qrcode
		$(GOGET) github.com/akavel/rsrc
		$(GOGET) gopkg.in/yaml.v3

build-apple-silicon:
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 $(GOBUILD) -o $(BINARY_NAME)-Silicon -v

build-apple-intel:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)-Intel -v

combine-universal:
	lipo -create -output $(BINARY_NAME) $(BINARY_NAME)-Intel $(BINARY_NAME)-Silicon

win-exe:
	$(GOBUILD) -o $(BINARY_NAME).exe -v

win-exe-prod:
	$(GOBUILD) -o $(BINARY_NAME).exe -v -ldflags -H=windowsgui

build-win32:
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v
other-win:
		GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ $(GOBUILD) -x -o $(BINARY_NAME)

win-manifest:
	rsrc -manifest ActionPadServer.manifest -arch amd64 -ico icons/icon.ico -o rsrc.syso
