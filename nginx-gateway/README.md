# nginx-gateway

Reverse proxy for backend services; serves API and docs.

## Build contexts
- **Dev Compose:** built from `../nginx-gateway` in `.dev/compose.yml`.
- **CI/CD:** built from `backend-services/nginx-gateway/Dockerfile` via GitHub Actions template.
- Build context must include `docs/specs/openapi.yaml` (generated via `scripts/generate-openapi.sh --all`).

## Runtime env vars (all required)
- `NGINX_GATEWAY_LISTEN_PORT` (default 8082)
- `NGINX_GATEWAY_BASE_UPSTREAM_HOST` / `PORT`
- `NGINX_GATEWAY_IMAGE_UPSTREAM_HOST` / `PORT`
- `NGINX_GATEWAY_LLM_UPSTREAM_HOST` / `PORT`
- `NGINX_GATEWAY_ACCOUNT_UPSTREAM_HOST` / `PORT`

## Nginx behavior (nginx.conf.template)
- Single location `^/v1/(base|album|llm|account)` proxies to mapped upstream pools.
- CORS + preflight handled; methods other than GET/POST/PUT/PATCH/DELETE/OPTIONS return 405.
- Security headers: HSTS, X-Frame-Options, X-Content-Type-Options, Referrer-Policy, Permissions-Policy; CSP for `/docs`.
- Limits: `limit_req_zone 20r/s` (burst 40, nodelay), `limit_conn` per IP 40, both return 429.
- Timeouts: connect 5s; send/read 30s; `client_max_body_size 50m`; header buffers tuned for uploads.

## Local dev
```
cd backend-services
docker compose -f .dev/compose.yml build nginx-gateway
docker compose -f .dev/compose.yml up -d nginx-gateway
```
Docs: http://localhost:8082/docs/

## Kubernetes (platform-infra)
- Deployment: `k8s/services/backend-services/backend-nginx-gateway-deployment.yml`.
- Apply & restart:
```
kubectl apply -f k8s/services/backend-services/backend-nginx-gateway-deployment.yml
kubectl rollout restart deployment/backend-nginx-gateway
```
Ingress: `api.nzhussup.dev` routes to `backend-nginx-gateway` service (8082).

## Quick rate-limit test
```
URL=${URL:-http://localhost:8082/v1/base/work-experience}
N=${N:-120} CONC=${CONC:-20}
export URL
seq "$N" | xargs -n1 -P"$CONC" -I{} sh -c 'curl -s -o /dev/null -w \"%{http_code}\\n\" \"$URL\"' | sort | uniq -c
```
- Expect 429 when limits are hit; adjust `limit_req` / `limit_conn` in `nginx.conf.template` if needed.
