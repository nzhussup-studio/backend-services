FROM node:24-alpine AS theme-builder

RUN apk add --no-cache maven openjdk17-jre-headless

WORKDIR /theme

COPY theme/package.json theme/package-lock.json ./
RUN npm ci --ignore-scripts

COPY theme/ ./
RUN npm run postinstall \
    && npm run build-keycloak-theme

FROM quay.io/keycloak/keycloak:26.1 AS keycloak-builder

COPY --from=theme-builder /theme/dist_keycloak/nzhussup-keycloak-theme.jar /opt/keycloak/providers/nzhussup-keycloak-theme.jar
RUN /opt/keycloak/bin/kc.sh build

FROM quay.io/keycloak/keycloak:26.1

USER root

COPY --from=keycloak-builder /opt/keycloak/ /opt/keycloak/

COPY --chown=keycloak:keycloak backend-auth-realm.template.json /opt/keycloak/conf/backend-auth-realm.template.json
COPY --chown=keycloak:keycloak --chmod=0755 docker-entrypoint.sh /opt/keycloak/bin/docker-entrypoint.sh

RUN mkdir -p /opt/keycloak/data/import \
    && chown -R keycloak:keycloak /opt/keycloak/data/import

USER keycloak

ENTRYPOINT ["/opt/keycloak/bin/docker-entrypoint.sh"]
CMD ["start", "--optimized", "--import-realm"]
