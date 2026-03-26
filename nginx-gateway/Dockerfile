FROM nginx:1.27-alpine

RUN apk add --no-cache gettext

COPY nginx-gateway/nginx.conf.template /etc/nginx/templates/nginx.conf.template
COPY nginx-gateway/docker-entrypoint.sh /docker-entrypoint.sh
COPY nginx-gateway/docs/index.html /usr/share/nginx/html/docs/index.html
COPY nginx-gateway/docs/favicon.svg /usr/share/nginx/html/docs/favicon.svg
COPY openapi.yaml /usr/share/nginx/html/docs/specs/openapi.yaml

# Ensure nginx (runs as non-root) can read the OpenAPI spec
RUN chmod 0644 /usr/share/nginx/html/docs/specs/openapi.yaml

RUN chmod 0755 /docker-entrypoint.sh

ENTRYPOINT ["/docker-entrypoint.sh"]
