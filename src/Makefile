build-windows:
	~/go/bin/templ generate
	GOOS=windows GOARCH=amd64 go build -o ./bin/itu-pickle.exe ./cmd

build-linux:
	~/go/bin/templ generate
	GOOS=linux GOARCH=amd64 go build -o ./bin/itu-pickle-linux ./cmd

build-macos:
	~/go/bin/templ generate
	GOOS=darwin GOARCH=amd64 go build -o ./bin/itu-pickle-macos ./cmd

build-arm:
	~/go/bin/templ generate
	GOOS=linux GOARCH=arm GOARM=6 go build -o ./bin/itu-pickle-arm ./cmd

build-all: build-windows build-linux build-macos build-arm
