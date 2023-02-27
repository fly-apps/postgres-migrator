ARG PG_VERSION=15.2

FROM golang:1.19 as flyutil

WORKDIR /go/src/github.com/fly-examples/fly-postgres
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -v -o /fly/bin/migrate ./cmd/migrate

FROM postgres:${PG_VERSION}

LABEL fly.app_role=postgres-migrator

RUN apt-get update && apt-get install --no-install-recommends -y \
    bash dnsutils \
    && apt autoremove -y

COPY --from=0 /fly/bin/* /usr/local/bin
