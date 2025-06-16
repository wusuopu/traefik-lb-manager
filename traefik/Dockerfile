FROM traefik:v3.4.0

COPY ./etc/traefik.yml /etc/traefik/traefik.yml
COPY ./plugins /app/plugins-local/src/github.com/wusuopu

WORKDIR /app/

ENV LB_MANAGER_LOG_FILE=/tmp/lb_manager_provider.log \
    LB_MANAGER_BASE_ENDPOINT= \
    LB_MANAGER_QUERY_PARAMS=
