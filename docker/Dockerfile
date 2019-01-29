
FROM ubuntu:latest

COPY sparks /usr/local/bin

RUN apt-get update

RUN sparks get -d --linux -v

ENTRYPOINT ["sparks"]
