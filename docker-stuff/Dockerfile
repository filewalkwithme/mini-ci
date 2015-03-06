FROM ubuntu

MAINTAINER Maicon Costa "maiconscosta@gmail.com"

WORKDIR /usr/local

RUN apt-get install -y wget git
RUN wget https://storage.googleapis.com/golang/go1.4.2.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.4.2.linux-amd64.tar.gz

ENV PATH /usr/local/go/bin:$PATH
ENV GOPATH /home/docker/go

RUN mkdir --parents /home/docker/go

ADD scripts/startup.sh /home/docker/scripts/startup.sh
ADD scripts/test.sh /home/docker/scripts/test.sh

CMD ["/home/docker/scripts/startup.sh"]
