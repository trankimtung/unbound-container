FROM debian:trixie AS unbound-build

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

FROM debian:trixie-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
        ca-certificates \
        libssl3 \
        libexpat1 \
        libevent-2.1-7 \
        libnghttp2-14 \
        libprotobuf-c1 \
        libsodium23 \
        ldnsutils \
        gettext-base \
    && rm -rf /var/lib/apt/lists/*

ENV UNBOUND_PUID=1000 \
    UNBOUND_GUID=1000 \
    UNBOUND_PORT=53

RUN groupadd -g ${UNBOUND_GUID} _unbound \
    && useradd -u ${UNBOUND_PUID} -g _unbound -s /usr/sbin/nologin _unbound

COPY --from=unbound-build /opt/unbound /opt/unbound

RUN mkdir -p /opt/unbound/etc/unbound/var \
    && chown -R _unbound:_unbound /opt/unbound/etc/unbound/var \
    && echo /opt/unbound/lib > /etc/ld.so.conf.d/unbound.conf \
    && ldconfig

COPY unbound.conf /opt/unbound/etc/unbound/unbound.conf.template
COPY entrypoint.sh /entrypoint.sh

ENV PATH=/opt/unbound/sbin:"$PATH"

LABEL org.opencontainers.image.version="1.24.2" \
      org.opencontainers.image.description="Unbound is a validating, recursive, and caching DNS resolver." \
      org.opencontainers.image.licenses="BSD-3-Clause" \
      org.opencontainers.image.source="https://github.com/trankimtung/unbound-container"

EXPOSE 53/tcp 53/udp

ENTRYPOINT ["/entrypoint.sh"]
