current_dir = $(shell pwd)
GO=go
# Give path of your go executable
#GO=/usr/lib/go-1.10/bin/go
# Path to your rsrc executable (see README.md)
RSRC=$(HOME)/go/bin/rsrc
MINGW_LIB?=$(HOME)/mingw-w64/current/lib

docker: clean
	mkdir -p bin
	-docker rm naksu-build
	docker build -t naksu-build-img -f Dockerfile.build .
	docker create --name naksu-build naksu-build-img
	docker cp naksu-build:/app/bin/ .

all: windows linux

windows: naksu.exe

linux: naksu

src/naksu.syso: res/windows/*
	$(RSRC) -arch="amd64" -ico="res/windows/naksu.ico" -o src/naksu.syso

naksu.exe: src/*
	cd src; GOPATH=$(current_dir)/ GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ CGO_LDFLAGS="-L$(MINGW_LIB)" $(GO) build -o ../bin/naksu.exe

naksu: src/*
	GOPATH=$(current_dir)/ GOARCH=amd64 $(GO) build -o bin/naksu src/naksu.go

naksu_packages: all
	rm -f naksu_linux_amd64.zip
	zip -j naksu_linux_amd64 bin/naksu
	rm -f naksu_windows_amd64.zip
	zip -j naksu_windows_amd64 bin/naksu.exe

update_libs:
	rm -fR src/github.com/
	rm -fR src/golang.org/
	#GOPATH=$(current_dir)/ go get github.com/gorilla/context
	GOPATH=$(current_dir)/ $(GO) get github.com/blang/semver
	GOPATH=$(current_dir)/ $(GO) get github.com/rhysd/go-github-selfupdate/selfupdate
	GOPATH=$(current_dir)/ $(GO) get github.com/andlabs/ui
	GOPATH=$(current_dir)/ $(GO) get github.com/kardianos/osext
	GOOS=windows GOPATH=$(current_dir)/ $(GO) get github.com/StackExchange/wmi
	GOPATH=$(current_dir)/ $(GO) get golang.org/x/text/encoding/charmap

clean:
	rm -f bin/naksu bin/naksu.exe

phony_get-server:
	VAGRANTPATH=phony-scripts/vagrant VBOXMANAGEPATH=phony-scripts/VBoxManage bin/get-server

phony_start-server:
	VAGRANTPATH=phony-scripts/vagrant VBOXMANAGEPATH=phony-scripts/VBoxManage bin/start-server
