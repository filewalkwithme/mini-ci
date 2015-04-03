FROM ubuntu

MAINTAINER Maicon Costa "maiconscosta@gmail.com"

WORKDIR /usr/local

#update
RUN apt-get update

#install go
RUN apt-get install -y wget git
RUN wget https://storage.googleapis.com/golang/go1.4.2.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.4.2.linux-amd64.tar.gz

#set environment variables for Go
ENV PATH /usr/local/go/bin:$PATH
ENV GOPATH /home/docker/go
RUN mkdir --parents /home/docker/go

#install postgres
RUN apt-get -y install postgresql postgresql-contrib

ADD scripts/startup.sh /home/docker/scripts/startup.sh

ADD keys/id_rsa /root/.ssh/id_rsa
ADD keys/id_rsa.pub /root/.ssh/id_rsa.pub
RUN chmod 600 -R /root/.ssh/

ADD minideploy/minideploy /usr/bin/minideploy

CMD ["/home/docker/scripts/startup.sh"]
