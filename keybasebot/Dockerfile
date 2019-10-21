ARG BUILD_FROM
FROM $BUILD_FROM 

MAINTAINER Prentice Wongvibulsin <me@prenticew.com> 

ENV KEYBASE_ALLOW_ROOT=1
ENV LANG C.UTF-8

RUN echo "http://dl-3.alpinelinux.org/alpine/edge/testing" >> /etc/apk/repositories
RUN apk upgrade --update-cache --available && apk add keybase-client

COPY keybasebot-armv7 /bin/keybasebot

RUN keybase ctl init
CMD ["/bin/keybasebot"]
