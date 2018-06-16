current_dir = $(shell pwd)
GO=go
# Give path of your go executable
#GO=/usr/lib/go-1.10/bin/go

all: windows linux

windows: naksu.exe

linux: naksu

naksu.exe: src/*
	GOPATH=$(current_dir)/ GOOS=windows GOARCH=amd64 $(GO) build -o bin/naksu.exe src/naksu.go

naksu: src/*
	GOPATH=$(current_dir)/ GOARCH=amd64 $(GO) build -o bin/naksu src/naksu.go

naksu_packages: all
	rm -f naksu_linux_amd64.zip
	zip -j naksu_linux_amd64 bin/naksu
	rm -f naksu_windows_amd64.zip
	zip -j naksu_windows_amd64 bin/naksu.exe

update_libs:
	rm -fR src/github.com/
	#GOPATH=$(current_dir)/ go get github.com/gorilla/context
	GOPATH=$(current_dir)/ $(GO) get github.com/blang/semver
	GOPATH=$(current_dir)/ $(GO) get github.com/rhysd/go-github-selfupdate/selfupdate

phony_get-server:
	VAGRANTPATH=phony-scripts/vagrant VBOXMANAGEPATH=phony-scripts/VBoxManage bin/get-server

phony_start-server:
	VAGRANTPATH=phony-scripts/vagrant VBOXMANAGEPATH=phony-scripts/VBoxManage bin/start-server
