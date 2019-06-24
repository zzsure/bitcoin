MAIN_PKG:=bitcoin
MAIN_PREFIX=$(dir $(MAIN_PKG))
MAIN=$(subst $(MAIN_PREFIX), , $(MAIN_PKG))
BIN=$(strip $(MAIN))

export GOPATH=$(shell pwd)/../../../../
export AIBEE_KUBERNETES_IDC=suzhou

build:
	go build -tags=jsoniter -x -o run/$(BIN) gitlab.azbit.cn/web/$(MAIN_PKG)

run: build
	cd run && ./$(BIN) $(ARG)

init:
	cd run && TARGET='run' ARG='init' docker-compose run --rm bitcoin-devel

docker-build:
	cd run && \
	cp $(BIN) ../build/ && \
	cd ../build && \
	docker build -t zzsure/bitcoin:$(TAG) . 
	#&& \
	#docker push zzsure/bitcoin:$(TAG)

.PHONY: build
