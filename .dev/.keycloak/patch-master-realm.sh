#!/bin/sh

set -eu

KEYCLOAK_URL="http://keycloak:8080"

echo "Waiting for Keycloak..."
until /opt/keycloak/bin/kcadm.sh config credentials \
  --server "${KEYCLOAK_URL}" \
  --realm master \
  --user "${KEYCLOAK_ADMIN_USERNAME}" \
  --password "${KEYCLOAK_ADMIN_PASSWORD}" >/dev/null 2>&1; do
  sleep 2
done

echo "Patching master realm SSL policy..."
/opt/keycloak/bin/kcadm.sh update realms/master -s sslRequired=NONE >/dev/null
echo "Master realm patched."
