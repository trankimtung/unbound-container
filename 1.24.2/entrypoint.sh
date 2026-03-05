#!/bin/sh
set -e

ANCHOR=/opt/unbound/etc/unbound/var/root.key
CONF=/opt/unbound/etc/unbound/unbound.conf
CONF_TEMPLATE=/opt/unbound/etc/unbound/unbound.conf.template

UNBOUND_PORT="${UNBOUND_PORT:-53}"

# Initialize or refresh the DNSSEC root trust anchor (RFC 5011).
# unbound-anchor exits non-zero if it updated the key (not an error).
# The var/ directory must be owned by the unbound user so that unbound
# can write the temp file used for atomic key updates after dropping privileges.
unbound-anchor -a "${ANCHOR}" || true

# Render config from template substituting environment variables.
UNBOUND_PORT="${UNBOUND_PORT}" envsubst '${UNBOUND_PORT}' < "${CONF_TEMPLATE}" > "${CONF}"

exec unbound -d -c "${CONF}"
