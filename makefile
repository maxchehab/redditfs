path := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

all: src/main.go
	cd src; go build -o ../dist/redditfs.exe
	redditfs ls

build: src/main.go 	
	cd src; go build -o ../dist/redditfs.exe

clean: 
	$(RM) ./dist/redditfs.exe

install:
	go get -u golang.org/x/sys/...
	go get github.com/fsnotify/fsnotify
	go get github.com/maxchehab/geddit
	go get golang.org/x/crypto/ssh
	go get github.com/ryanuber/columnize
	go get gopkg.in/AlecAivazis/survey.v1
	go get github.com/ttacon/chalk

run: dist/redditfs.exe
	./dist/redditfs.exe

link: dist/redditfs.exe
	sudo ln -s $(path)/dist/redditfs.exe /bin/redditfs

test: 
	cd src; go test