# api-docs

Static Swagger UI for the unified backend OpenAPI document. API routing, authentication,
CORS, and rate limiting are handled by Envoy Gateway; this container only serves `/docs/`.

Regenerate the bundled specification with:

```sh
./scripts/generate-openapi.sh --all
```

Run the complete local stack and open <http://localhost:8082/docs/>. In Kubernetes,
Envoy Gateway routes `/docs` to the `backend-api-docs` service.
