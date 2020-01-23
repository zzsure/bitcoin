MAIN_PKG:=bitcoin
MAIN_PREFIX=$(dir $(MAIN_PKG))
MAIN=$(subst $(MAIN_PREFIX), , $(MAIN_PKG))
BIN=$(strip $(MAIN))

export GOPATH=$(shell pwd)/../../../../
export AZBIT_KUBERNETES_IDC=suzhou
export GITTAG=$(shell git describe --tags `git rev-list --tags --max-count=1`)
export GITHASH=$(shell git rev-list HEAD -n 1 | cut -c 1-)
export GITBRANCH=$(shell git symbolic-ref --short -q HEAD)

build:
	go build -tags=jsoniter -x -o run/$(BIN) . 

run: build
	cd run && ./$(BIN) $(ARG)

init:
	cd run && TARGET='run' ARG='init' docker-compose run --rm bitcoin-devel

docker-build:
	docker build . -t zzsure/bitcoin:$(GITTAG) && \
	docker push zzsure/bitcoin:$(GITTAG)

.PHONY: build
