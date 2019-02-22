.PHONY: build_linux32 build_linux64 build_win32 build_win64

build: build_linux32 build_linux64 build_win32 build_win64

build_win64:
	@GOOS=windows GOARCH=amd64 go build -o bin/goat.exe . && ls -sh bin/goat.exe

build_win32:
	@GOOS=windows GOARCH=386 go build -o bin/goat_x32.exe . && ls -sh bin/goat_x32.exe

build_linux64:
	@GOOS=linux GOARCH=amd64 go build -o bin/goat . && ls -sh bin/goat

build_linux32:
	@GOOS=linux GOARCH=386 go build -o bin/goat_x32 . && ls -sh bin/goat_x32
