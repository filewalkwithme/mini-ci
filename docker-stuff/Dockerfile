FROM ubuntu

MAINTAINER Maicon Costa "maiconscosta@gmail.com"

RUN apt-get install -y wget git

WORKDIR /usr/local

RUN wget https://storage.googleapis.com/golang/go1.4.2.linux-amd64.tar.gz

RUN tar -C /usr/local -xzf go1.4.2.linux-amd64.tar.gz

ENV PATH /usr/local/go/bin:$PATH

ENV GOPATH /home/docker/go


RUN mkdir --parents /home/docker/go

VOLUME /home/docker/scripts

CMD ["/home/docker/scripts/startup.sh"]

