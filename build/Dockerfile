FROM golang:1.10-alpine 
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai/etc/localtime
WORKDIR /opt/bitcoin

RUN mkdir -p /opt/bitcoin
COPY config.toml /opt/bitcoin/config.toml
COPY bitcoin /opt/bitcoin/bitcoin
RUN chmod +x /opt/bitcoin/bitcoin

ENTRYPOINT ["/opt/bitcoin/bitcoin"]
