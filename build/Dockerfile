FROM alpine:3.7
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai/etc/localtime
WORKDIR /opt/bitcoin
COPY config.toml /opt/bitcoin/config.toml
COPY bitcoin /opt/bitcoin/bitcoin
RUN chmod +x /opt/bitcoin/bitcoin
ENTRYPOINT ["/opt/bitcoin/bitcoin"]
