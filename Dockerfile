FROM golang:1.13.3 as build
WORKDIR /root/code/zzsure/bitcoin
COPY . .
RUN go build -mod=vendor -ldflags "-X 'main.goversion=$(go version)'" -o bitcoin main.go

FROM ubuntu:xenial
RUN apt-get update
RUN apt-get install tzdata
RUN echo "Asia/Shanghai" > /etc/timezone
RUN rm -f /etc/localtime
RUN dpkg-reconfigure -f noninteractive tzdata
RUN apt-get install -y ca-certificates

WORKDIR /root/deploy/bitcoin
COPY --from=build /root/code/zzsure/bitcoin/bitcoin ./bitcoin
COPY --from=build /root/code/zzsure/bitcoin/run/config.toml ./config.toml
