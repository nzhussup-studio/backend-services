#!/usr/bin/env bash

set -euo pipefail

deployments="${1:-}"

if [[ -z "$deployments" ]]; then
  exit 0
fi

for deployment in $deployments; do
  kubectl rollout restart deployment/"$deployment"
done
