current_dir = $(shell pwd)
GO=go
# Give path of your go executable
#GO=/usr/lib/go-1.10/bin/go
# Path to your rsrc executable (see README.md)
RSRC=$(HOME)/go/bin/rsrc
MINGW_LIB?=$(HOME)/mingw-w64/current/lib

bin/gometalinter:
	curl https://raw.githubusercontent.com/alecthomas/gometalinter/master/scripts/install.sh | sh

checkstyle: bin/gometalinter
	-GOOS=linux GOARCH=amd64 CGO_ENABLED=1 ./bin/gometalinter --deadline=240s --vendor --checkstyle ./src/naksu/... > checkstyle-linux.xml
	-GOOS=windows GOARCH=amd64 CGO_ENABLED=1 ./bin/gometalinter --deadline=240s --vendor --checkstyle ./src/naksu/... > checkstyle-windows.xml

lint: bin/gometalinter
	./bin/gometalinter --deadline=240s --vendor ./src/naksu/...

docker: clean
	mkdir -p bin
	-docker rm naksu-build
	docker build -t naksu-build-img:latest -f Dockerfile.build .
	docker create --name naksu-build naksu-build-img
	docker cp naksu-build:/app/checkstyle-linux.xml .
	docker cp naksu-build:/app/checkstyle-windows.xml .
	docker cp naksu-build:/app/bin/naksu bin/naksu
	docker cp naksu-build:/app/bin/naksu.exe bin/naksu.exe
	docker cp naksu-build:/app/naksu_linux_amd64.zip .
	docker cp naksu-build:/app/naksu_windows_amd64.zip .

all: windows linux

windows: naksu.exe

linux: naksu

src/naksu.syso: res/windows/*
	$(RSRC) -arch="amd64" -ico="res/windows/naksu.ico" -o src/naksu.syso

naksu.exe: src/*
	GOPATH=$(current_dir)/ GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ CGO_LDFLAGS="-L$(MINGW_LIB)" $(GO) build -o bin/naksu.exe naksu

naksu: src/*
	GOPATH=$(current_dir)/ GOOS=linux GOARCH=amd64 CGO_ENABLED=1 $(GO) build -o bin/naksu naksu

naksu_packages: all
	rm -f naksu_linux_amd64.zip
	zip -j naksu_linux_amd64 bin/naksu
	rm -f naksu_windows_amd64.zip
	zip -j naksu_windows_amd64 bin/naksu.exe

update_libs: clean
	cd src/naksu && dep ensure

clean:
	rm -f bin/naksu bin/naksu.exe

phony_get-server:
	VAGRANTPATH=phony-scripts/vagrant VBOXMANAGEPATH=phony-scripts/VBoxManage bin/get-server

phony_start-server:
	VAGRANTPATH=phony-scripts/vagrant VBOXMANAGEPATH=phony-scripts/VBoxManage bin/start-server
