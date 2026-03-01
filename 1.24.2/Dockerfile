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
    && groupadd -r -g 1000 unbound \
    && useradd -r -u 1000 -g unbound -s /usr/sbin/nologin unbound \
    && ./configure -h \
    && ./configure \
        --prefix=/opt/unbound \
        --with-username=unbound \
        --with-libevent \
        --with-ssl \
        --with-pthreads \
        --enable-subnet \
        --enable-dnstap \
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
    && rm -rf /var/lib/apt/lists/*

RUN groupadd -r -g 1000 unbound \
    && useradd -r -u 1000 -g unbound -s /usr/sbin/nologin unbound

RUN mkdir -p /opt/unbound

COPY --from=unbound-build /opt/unbound /opt/unbound
RUN echo /opt/unbound/lib > /etc/ld.so.conf.d/unbound.conf && ldconfig

COPY unbound.conf /opt/unbound/etc/unbound/unbound.conf
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENV PATH=/opt/unbound/sbin:"$PATH"

EXPOSE 53/tcp 53/udp

ENTRYPOINT ["/entrypoint.sh"]