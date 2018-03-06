path := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

all: main.go
	go build

build: main.go
	go build

clean: 
	$(RM) redditfs

install:
	go get -u golang.org/x/sys/...
	go get github.com/fsnotify/fsnotify
	go get github.com/maxchehab/geddit
	go get golang.org/x/crypto/ssh
	go get github.com/ryanuber/columnize
	go get gopkg.in/AlecAivazis/survey.v1
	go get github.com/ttacon/chalk
	go get github.com/monochromegane/go-gitignore

run: dist/redditfs.exe
	./redditfs

link: dist/redditfs.exe
	sudo ln -s $(path)redditfs /bin/redditfs