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

assert_status 200 /v1/base/health
assert_status 200 /v1/account/health
assert_status 200 '/v1/album?type=public'
assert_status 200 /docs/
assert_status 401 /v1/account
assert_status 401 /v1/base/cv-generator-preference
assert_status 401 '/v1/album?type=all'
assert_status 401 '/v1/album?type=private'
assert_status 401 '/v1/album?type=semi-public'

echo "Gateway smoke tests passed."
