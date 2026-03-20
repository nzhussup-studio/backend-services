#!/bin/sh

set -eu

KEYCLOAK_URL="http://keycloak:8080"
ACCOUNT_SERVICE_ADMIN_CLIENT_ID="${ACCOUNT_SERVICE_ADMIN_CLIENT_ID:-account-service-admin}"
ACCOUNT_SERVICE_ADMIN_CLIENT_SECRET="${ACCOUNT_SERVICE_ADMIN_CLIENT_SECRET:-account-service-dev-secret}"
TARGET_REALM="${KEYCLOAK_REALM_NAME:-backend-auth-dev}"

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

echo "Ensuring dedicated account-service admin client exists in master..."
CLIENT_ID=$(
  /opt/keycloak/bin/kcadm.sh get clients -r master -q clientId="${ACCOUNT_SERVICE_ADMIN_CLIENT_ID}" --fields id,clientId \
    | sed -n 's/.*"id"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' \
    | head -n 1
)

if [ -z "${CLIENT_ID}" ]; then
  /opt/keycloak/bin/kcadm.sh create clients -r master \
    -s clientId="${ACCOUNT_SERVICE_ADMIN_CLIENT_ID}" \
    -s enabled=true \
    -s protocol=openid-connect \
    -s publicClient=false \
    -s serviceAccountsEnabled=true \
    -s standardFlowEnabled=false \
    -s directAccessGrantsEnabled=false \
    -s bearerOnly=false \
    -s secret="${ACCOUNT_SERVICE_ADMIN_CLIENT_SECRET}" >/dev/null

  CLIENT_ID=$(
    /opt/keycloak/bin/kcadm.sh get clients -r master -q clientId="${ACCOUNT_SERVICE_ADMIN_CLIENT_ID}" --fields id,clientId \
      | sed -n 's/.*"id"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' \
      | head -n 1
  )
else
  /opt/keycloak/bin/kcadm.sh update "clients/${CLIENT_ID}" -r master \
    -s enabled=true \
    -s publicClient=false \
    -s serviceAccountsEnabled=true \
    -s standardFlowEnabled=false \
    -s directAccessGrantsEnabled=false \
    -s bearerOnly=false \
    -s secret="${ACCOUNT_SERVICE_ADMIN_CLIENT_SECRET}" >/dev/null
fi

SERVICE_ACCOUNT_USERNAME=$( \
  /opt/keycloak/bin/kcadm.sh get "clients/${CLIENT_ID}/service-account-user" -r master \
    | sed -n 's/.*"username"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' \
)

if [ -z "${SERVICE_ACCOUNT_USERNAME}" ]; then
  echo "Failed to resolve service account for ${ACCOUNT_SERVICE_ADMIN_CLIENT_ID}" >&2
  exit 1
fi

echo "Granting ${TARGET_REALM}-realm roles to ${SERVICE_ACCOUNT_USERNAME} in master..."
/opt/keycloak/bin/kcadm.sh add-roles \
  -r master \
  --uusername "${SERVICE_ACCOUNT_USERNAME}" \
  --cclientid "${TARGET_REALM}-realm" \
  --rolename create-client \
  --rolename impersonation \
  --rolename manage-authorization \
  --rolename manage-clients \
  --rolename manage-events \
  --rolename manage-identity-providers \
  --rolename manage-realm \
  --rolename manage-users \
  --rolename query-clients \
  --rolename query-groups \
  --rolename query-realms \
  --rolename query-users \
  --rolename view-authorization \
  --rolename view-clients \
  --rolename view-events \
  --rolename view-identity-providers \
  --rolename view-realm \
  --rolename view-users >/dev/null

echo "Master realm patched."
