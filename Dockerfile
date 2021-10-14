ARG PG_VERSION=13.4

FROM postgres:${PG_VERSION}

RUN apt-get update && apt-get install --no-install-recommends -y \
    ca-certificates curl bash dnsutils \
    && apt autoremove -y

ADD ./migrate /usr/local/bin

