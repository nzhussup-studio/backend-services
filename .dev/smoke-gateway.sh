#!/bin/sh

set -eu

gateway_url="${API_GATEWAY_URL:-http://localhost:8082}"

assert_status() {
  expected="$1"
  path="$2"
  actual="$(curl --silent --show-error --output /dev/null --write-out '%{http_code}' "${gateway_url}${path}")"

  if [ "${actual}" != "${expected}" ]; then
    echo "Expected ${expected} for ${path}, received ${actual}" >&2
    exit 1
  fi

  echo "${actual} ${path}"
}

assert_preflight() {
  path="$1"
  headers="$(curl --silent --show-error --dump-header - --output /dev/null \
    --request OPTIONS "${gateway_url}${path}" \
    --header 'Origin: http://localhost:5173' \
    --header 'Access-Control-Request-Method: GET' \
    --header 'Access-Control-Request-Headers: authorization')"

  echo "${headers}" | grep -qi '^HTTP/.* 200'
  echo "${headers}" | grep -qi '^access-control-allow-origin: http://localhost:5173'
  echo "200 OPTIONS ${path}"
}

assert_status 200 /v1/albums
assert_status 200 /docs/
assert_status 401 /v1/me/account
assert_status 401 /v1/admin/cv-generator-preferences
assert_status 401 /v1/admin/albums
assert_status 401 /v1/admin/llm/configuration
assert_preflight /v1/admin/skills

echo "Gateway smoke tests passed."
