current_dir = $(shell pwd)

all: windows linux

windows: install.exe start-server.exe

linux: install start-server

install.exe: src/install.go src/mebroutines/*
	GOPATH=$(current_dir)/ GOOS=windows GOARCH=386 go build -o bin/install.exe src/install.go

install: src/install.go src/mebroutines/*
	GOPATH=$(current_dir)/ GOARCH=386 go build -o bin/install src/install.go

start-server.exe: src/start-server.go src/mebroutines/*
	GOPATH=$(current_dir)/ GOOS=windows GOARCH=386 go build -o bin/start-server.exe src/start-server.go

start-server: src/start-server.go src/mebroutines/*
	GOPATH=$(current_dir)/ GOARCH=386 go build -o bin/start-server src/start-server.go

phony_install:
	VAGRANTPATH=phony-scripts/vagrant VBOXMANAGEPATH=phony-scripts/VBoxManage bin/install

phony_start-server:
	VAGRANTPATH=phony-scripts/vagrant VBOXMANAGEPATH=phony-scripts/VBoxManage bin/install
