current_dir = $(shell pwd)

all: windows linux

windows: get-server.exe start-server.exe

linux: get-server start-server

get-server.exe: src/get-server.go src/mebroutines/*
	GOPATH=$(current_dir)/ GOOS=windows GOARCH=386 go build -o bin/get-server.exe src/get-server.go

get-server: src/get-server.go src/mebroutines/*
	GOPATH=$(current_dir)/ GOARCH=386 go build -o bin/get-server src/get-server.go

start-server.exe: src/start-server.go src/mebroutines/*
	GOPATH=$(current_dir)/ GOOS=windows GOARCH=386 go build -o bin/start-server.exe src/start-server.go

start-server: src/start-server.go src/mebroutines/*
	GOPATH=$(current_dir)/ GOARCH=386 go build -o bin/start-server src/start-server.go

phony_get-server:
	VAGRANTPATH=phony-scripts/vagrant VBOXMANAGEPATH=phony-scripts/VBoxManage bin/get-server

phony_start-server:
	VAGRANTPATH=phony-scripts/vagrant VBOXMANAGEPATH=phony-scripts/VBoxManage bin/start-server
