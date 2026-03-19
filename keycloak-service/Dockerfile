FROM quay.io/keycloak/keycloak:26.1

USER root

COPY --chown=keycloak:keycloak backend-auth-realm.template.json /opt/keycloak/conf/backend-auth-realm.template.json
COPY --chown=keycloak:keycloak --chmod=0755 docker-entrypoint.sh /opt/keycloak/bin/docker-entrypoint.sh

RUN mkdir -p /opt/keycloak/data/import \
    && chown -R keycloak:keycloak /opt/keycloak/data/import

USER keycloak

ENTRYPOINT ["/opt/keycloak/bin/docker-entrypoint.sh"]
CMD ["start", "--import-realm"]
