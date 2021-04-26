FROM ubuntu:latest
MAINTAINER melvinto@gmail.com

COPY bin/promhub /usr/bin/promhub
RUN mkdir /etc/promhub

