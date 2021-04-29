FROM ubuntu:latest
MAINTAINER melvinto@gmail.com

RUN apt update
RUN apt install -y ca-certificates
RUN apt clean


COPY bin/promhub /usr/bin/promhub
RUN mkdir /etc/promhub

CMD /usr/bin/promhub -c /etc/promhub/config.yml
