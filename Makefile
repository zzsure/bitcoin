MAIN_PKG:=bitcoin
MAIN_PREFIX=$(dir $(MAIN_PKG))
MAIN=$(subst $(MAIN_PREFIX), , $(MAIN_PKG))
BIN=$(strip $(MAIN))

export GOPATH=$(shell pwd)/../../../../
export AIBEE_KUBERNETES_IDC=suzhou

build:
	go build -tags=jsoniter -x -o run/$(BIN) gitlab.azbit.cn/web/$(MAIN_PKG)

dev:
	glr -main gitlab.azbit.cn/web/$(MAIN_PKG) -wd run -delay 2000 -args $(ARG)

run: build
	cd run && ./$(BIN) $(ARG)

init:
	cd run && TARGET='run' ARG='init' docker-compose run --rm HelloWorld-devel

docker-build:
	cd run && \
	TARGET='build' docker-compose run --rm HelloWorld-devel && cp $(BIN) ../build/ && \
	cd ../build && \
	docker build -t hub.baidubce.com/azbit/HelloWorld:$(TAG) . && \
	docker push hub.baidubce.com/azbit/HelloWorld:$(TAG)

.PHONY: build
