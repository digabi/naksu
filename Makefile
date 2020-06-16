current_dir = $(shell pwd)
GO=go
# Give path of your go executable
# GO=/usr/lib/go-1.10/bin/go
# Path to your rsrc executable (see README.md)
RSRC=$(HOME)/go/bin/rsrc
TESTS=naksu/mebroutines naksu/mebroutines/backup naksu naksu/box naksu/host naksu/boxversion naksu/network

bin/gometalinter:
	curl https://raw.githubusercontent.com/alecthomas/gometalinter/master/scripts/install.sh | sh

bin/go2xunit:
	GOPATH=$(current_dir)/ go get github.com/tebeka/go2xunit

checkstyle: bin/gometalinter
	-GOOS=linux GOARCH=amd64 CGO_ENABLED=1 ./bin/gometalinter --deadline=240s --vendor --checkstyle ./src/naksu/... > checkstyle-linux.xml
	-GOOS=windows GOARCH=amd64 CGO_ENABLED=1 ./bin/gometalinter --deadline=240s --vendor --checkstyle ./src/naksu/... > checkstyle-windows.xml

lint: bin/gometalinter
	./bin/gometalinter --deadline=240s --vendor ./src/naksu/...

ci-test: bin/go2xunit
	2>&1 GOPATH=$(current_dir)/ go test -v $(TESTS) | ./bin/go2xunit -output tests.xml

test:
	GOPATH=$(current_dir)/ go test $(TESTS)

docker: clean
	mkdir -p bin
	-docker rm naksu-build
	docker build -t naksu-build-img:latest -f Dockerfile.build .
	docker run -w /app --name naksu-build naksu-build-img:latest make ci-test
	docker cp naksu-build:/app/checkstyle-linux.xml .
	docker cp naksu-build:/app/checkstyle-windows.xml .
	docker cp naksu-build:/app/tests.xml .
	docker cp naksu-build:/app/bin/naksu bin/naksu
	docker cp naksu-build:/app/bin/naksu.exe bin/naksu.exe
	docker cp naksu-build:/app/naksu_linux_amd64.zip .
	docker cp naksu-build:/app/naksu_windows_amd64.zip .

all: test windows linux

windows: naksu.exe

linux: naksu

mac: naksu-darwin

src/naksu.syso: res/windows/*
	$(RSRC) -arch="amd64" -ico="res/windows/naksu.ico" -o src/naksu.syso

naksu.exe: src/*
	GOPATH=$(current_dir)/ GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ $(GO) build -o bin/naksu.exe naksu

naksu: src/*
	GOPATH=$(current_dir)/ GOOS=linux GOARCH=amd64 CGO_ENABLED=1 $(GO) build -o bin/naksu naksu

naksu-darwin: src/*
	GOPATH=$(current_dir)/ GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 $(GO) build -o bin/naksu-darwin naksu

naksu_packages: all
	rm -f naksu_linux_amd64.zip
	zip -j naksu_linux_amd64 bin/naksu
	rm -f naksu_windows_amd64.zip
	zip -j naksu_windows_amd64 bin/naksu.exe

update_libs: clean
	cd src/naksu && GOPATH=$(current_dir) dep ensure

clean:
	rm -f bin/naksu bin/naksu.exe
	rm -f tests.xml

phony_get-server:
	VAGRANTPATH=phony-scripts/vagrant VBOXMANAGEPATH=phony-scripts/VBoxManage bin/get-server

phony_start-server:
	VAGRANTPATH=phony-scripts/vagrant VBOXMANAGEPATH=phony-scripts/VBoxManage bin/start-server
