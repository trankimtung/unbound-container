#!/bin/sh
set -e

ANCHOR=/opt/unbound/etc/unbound/var/root.key

# Initialize or refresh the DNSSEC root trust anchor (RFC 5011).
# unbound-anchor exits non-zero if it updated the key (not an error).
# The var/ directory must be owned by the unbound user so that unbound
# can write the temp file used for atomic key updates after dropping privileges.
mkdir -p /opt/unbound/etc/unbound/var
unbound-anchor -a "${ANCHOR}" || true
chown -R unbound:unbound /opt/unbound/etc/unbound/var

exec unbound -d -c /opt/unbound/etc/unbound/unbound.conf
