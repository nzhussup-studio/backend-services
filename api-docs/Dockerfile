FROM nginx:1.29.2-alpine

COPY default.conf /etc/nginx/conf.d/default.conf
COPY docs/ /usr/share/nginx/html/docs/

EXPOSE 8080

HEALTHCHECK --interval=10s --timeout=5s --retries=5 \
  CMD wget --quiet --output-document=- http://127.0.0.1:8080/docs/ >/dev/null || exit 1
