VERSION=latest
build-docker: Dockerfile
	docker build -t hedykan/http-proxy:$(VERSION) .
	
push-docker: build-docker
	docker tag hedykan/http-proxy:$(VERSION) registry.cn-hangzhou.aliyuncs.com/hedykan/http-proxy:$(VERSION)
	docker push registry.cn-hangzhou.aliyuncs.com/hedykan/http-proxy:$(VERSION)

build-all: build-linux build-window build-docker

build-window:
	go build

build-linux:
	GOOS=linux go build

Dockerfile:
	goctl docker -go main.go

clean:
	rm http-proxy* -f