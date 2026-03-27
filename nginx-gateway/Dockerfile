FROM nginx:1.27-alpine

RUN apk add --no-cache gettext

COPY ./nginx.conf.template /etc/nginx/templates/nginx.conf.template
COPY ./docker-entrypoint.sh /docker-entrypoint.sh
COPY ./docs/index.html /usr/share/nginx/html/docs/index.html
COPY ./docs/favicon.svg /usr/share/nginx/html/docs/favicon.svg
COPY ./docs/specs/openapi.yaml /usr/share/nginx/html/docs/specs/openapi.yaml

# Ensure nginx (runs as non-root) can read the OpenAPI spec
RUN chmod 0644 /usr/share/nginx/html/docs/specs/openapi.yaml

RUN chmod 0755 /docker-entrypoint.sh

ENTRYPOINT ["/docker-entrypoint.sh"]
