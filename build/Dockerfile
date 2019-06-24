FROM ubuntu:xenial
RUN apt-get update
RUN apt-get install tzdata
RUN echo "Asia/Shanghai" > /etc/timezone
RUN rm -f /etc/localtime
RUN dpkg-reconfigure -f noninteractive tzdata
RUN apt-get install -y ca-certificates
#RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai/etc/localtime
#RUN ln -fs /usr/share/zoneinfo/US/Pacific-New /etc/localtime && dpkg-reconfigure -f noninteractive tzdata
WORKDIR /opt/bitcoin

RUN mkdir -p /opt/bitcoin
COPY config.toml /opt/bitcoin/config.toml
COPY bitcoin /opt/bitcoin/bitcoin
RUN chmod +x /opt/bitcoin/bitcoin

ENTRYPOINT ["/opt/bitcoin/bitcoin"]
