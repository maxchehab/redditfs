all: src/main.go
	cd src; go build -o ../dist/redditfs.exe
	./dist/redditfs.exe

build: src/main.go 	
	cd src; go build -o ../dist/redditfs.exe

clean: 
	$(RM) ./dist/redditfs.exe

install:
	go get -u golang.org/x/sys/...
	go get github.com/fsnotify/fsnotify

run: dist/redditfs.exe
	./dist/redditfs.exe