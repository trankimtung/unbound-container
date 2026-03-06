FROM debian:trixie-slim AS unbound-build

ARG UNBOUND_VERSION=1.24.2
ARG UNBOUND_SHA256=44e7b53e008a6dcaec03032769a212b46ab5c23c105284aa05a4f3af78e59cdb
ARG UNBOUND_DOWNLOAD_URL=https://nlnetlabs.nl/downloads/unbound/unbound-${UNBOUND_VERSION}.tar.gz

RUN apt-get update && apt-get install -y --no-install-recommends \
        build-essential \
        curl \
        ca-certificates \
        libssl-dev \
        libexpat1-dev \
        libevent-dev \
        libnghttp2-dev \
        libprotobuf-c-dev \
        protobuf-c-compiler \
        libsodium-dev \
        bison \
        flex \
        autoconf \
        automake \
        pkg-config \
    && rm -rf /var/lib/apt/lists/*

RUN curl -fL -o unbound.tar.gz ${UNBOUND_DOWNLOAD_URL} \
    && echo "${UNBOUND_SHA256}  unbound.tar.gz" | sha256sum -c - \
    && tar -xzf unbound.tar.gz \
    && cd unbound-${UNBOUND_VERSION} \
    && ./configure \
        --prefix=/opt/unbound \
        --with-username=_unbound \
        --with-libevent \
        --with-ssl \
        --with-pthreads \
        --without-pyunbound \
        --without-pythonmodule \
        --enable-subnet \
        --enable-dnstap \
        --enable-dnscrypt \
        --enable-cachedb \
        --enable-ipsecmod \
        --enable-event-api \
        --enable-tfo-server \
        --enable-tfo-client \
        --enable-relro-now \
        --enable-year2038 \
    && make -j$(nproc) \
    && make install \
    && mv /opt/unbound/etc/unbound/unbound.conf /opt/unbound/etc/unbound/unbound.conf.example

FROM debian:trixie-slim AS lib-provider

RUN apt-get update && apt-get install -y --no-install-recommends \
        ca-certificates \
        libssl3 \
        libexpat1 \
        libevent-2.1-7 \
        libnghttp2-14 \
        libprotobuf-c1 \
        libsodium23 \
        ldnsutils \
    && rm -rf /var/lib/apt/lists/*

COPY --from=unbound-build /opt/unbound /opt/unbound

RUN groupadd -g 1000 _unbound \
    && useradd -u 1000 -g _unbound -s /usr/sbin/nologin -M _unbound \
    && mkdir -p /opt/unbound/etc/unbound/var \
    && chown -R _unbound:_unbound /opt/unbound/etc/unbound

# Collect required .so files into a flat directory, dereferencing symlinks.
RUN mkdir -p /runtime-libs \
    && find /usr/lib -maxdepth 2 \( \
        -name "libssl.so*" \
        -o -name "libcrypto.so*" \
        -o -name "libexpat.so*" \
        -o -name "libevent-2.1.so*" \
        -o -name "libnghttp2.so*" \
        -o -name "libprotobuf-c.so*" \
        -o -name "libsodium.so*" \
        -o -name "libldns.so*" \
    \) -exec cp -L {} /runtime-libs/ \;

FROM golang:1.25-trixie AS entrypoint-build

WORKDIR /build
COPY entrypoint/go.mod ./
COPY entrypoint/main.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /entrypoint .

FROM gcr.io/distroless/base-debian13

COPY --from=lib-provider /runtime-libs /usr/local/lib
COPY --from=lib-provider /etc/ssl/certs /etc/ssl/certs
COPY --from=lib-provider /etc/passwd /etc/passwd
COPY --from=lib-provider /etc/group /etc/group
COPY --from=lib-provider /opt/unbound /opt/unbound
COPY --from=entrypoint-build /entrypoint /entrypoint
COPY --from=lib-provider /usr/bin/drill /usr/bin/drill
COPY unbound.conf /opt/unbound/etc/unbound/unbound.conf.template

ENV PATH=/opt/unbound/sbin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin \
    LD_LIBRARY_PATH=/usr/local/lib:/opt/unbound/lib \
    UNBOUND_PORT=5335

USER _unbound

LABEL org.opencontainers.image.version="1.24.2" \
      org.opencontainers.image.description="Unbound is a validating, recursive, and caching DNS resolver." \
      org.opencontainers.image.licenses="BSD-3-Clause" \
      org.opencontainers.image.source="https://github.com/trankimtung/unbound-container"

EXPOSE 5335/tcp 5335/udp

ENTRYPOINT ["/entrypoint"]
