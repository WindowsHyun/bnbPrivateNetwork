# syntax=docker/dockerfile:1
FROM ubuntu:20.04
#RUN rm /bin/sh && ln -s /bin/bash /bin/sh
ARG DEBIAN_FRONTEND=noninteractive
ENV PATH=/usr/local/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
ENV TZ=Asia/Seoul
ENV LANG=C.UTF-8
RUN apt update && apt install -y --no-install-recommends vim net-tools curl openssh-server ca-certificates git software-properties-common tzdata unzip wget jq
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
RUN mkdir -p /workdir
RUN mkdir -p /localnet
RUN mkdir -p /bnbsmartchain 
WORKDIR /localnet

SHELL ["/bin/bash", "-c"]

RUN cd /bnbsmartchain  \
    && wget $(curl -s https://api.github.com/repos/bnb-chain/bsc/releases/latest |grep browser_ |grep geth_linux |cut -d\" -f4) \
    && mv geth_linux /usr/local/bin/geth 
RUN apt purge -y software-properties-common
RUN apt autoremove -y
RUN rm -rf /var/lib/apt/lists/*

ADD config.toml /bnbsmartchain/config.toml
ADD genesis.json /bnbsmartchain/genesis.json

RUN chmod +x /usr/local/bin/geth
RUN chmod +x /bnbsmartchain/config.toml
RUN chmod +x /bnbsmartchain/genesis.json

ADD sshd_config /sshd_config
COPY ./sshd_config /etc/ssh/sshd_config
RUN chmod 644 /etc/ssh/sshd_config
COPY ./start.sh /start.sh
RUN chmod +x /start.sh
RUN (echo 'localhost_root'; echo 'localhost_root') | passwd root
RUN update-rc.d ssh defaults
CMD /start.sh