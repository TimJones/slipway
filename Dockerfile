FROM golang:1.11 AS development
MAINTAINER Tim Jones <timniverse@gmail.com>
ARG UID=1000
ARG GID=1000
ARG USER=slipdev
RUN set -eux; \
    curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh; \
    groupadd --gid ${GID} ${USER}; \
    useradd --create-home --uid ${UID} --gid ${GID} ${USER}; \
    chown -R ${USER} ${GOPATH}
WORKDIR /go/src/github.com/timjones/slipway

FROM development AS build
COPY . ./
RUN set -eux; \
    dep ensure -v; \
    ./scripts/build;

FROM alpine:latest
COPY --from=build /go/src/github.com/timjones/slipway/bin/* /usr/local/bin/
ENTRYPOINT /usr/local/bin/slip
