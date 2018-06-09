current_dir = $(shell pwd)

all: windows linux

windows: naksu.exe

linux: naksu

naksu.exe: src/*
	GOPATH=$(current_dir)/ GOOS=windows GOARCH=386 go build -o bin/naksu.exe src/naksu.go

naksu: src/*
	GOPATH=$(current_dir)/ GOARCH=386 go build -o bin/naksu src/naksu.go

phony_get-server:
	VAGRANTPATH=phony-scripts/vagrant VBOXMANAGEPATH=phony-scripts/VBoxManage bin/get-server

phony_start-server:
	VAGRANTPATH=phony-scripts/vagrant VBOXMANAGEPATH=phony-scripts/VBoxManage bin/start-server
