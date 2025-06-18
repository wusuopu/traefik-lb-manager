# Stage 1 - frontend
FROM node:20-alpine as frontend

WORKDIR /app

COPY ./frontend/package.json ./frontend/yarn.lock /app/
RUN cd /app && yarn install && rm -rf /root/.cache /root/.npm /usr/local/share/.cache/yarn/ /tmp/yarn*

COPY ./frontend/ /app
RUN cd /app && yarn run build-only


# Stage 2 - backend
FROM golang:1.24.2-alpine as backend

ENV GOPROXY=https://goproxy.cn,direct \
    GO111MODULE=auto

WORKDIR /app/

COPY ./backend/go.mod ./backend/go.sum /app/

RUN go mod download

COPY ./backend/ /app
RUN rm -rf assets

COPY --from=frontend /app/dist /app/assets

RUN go build -ldflags "-s -w" -o dist/goose cmd/goose.go && \
    go build -tags=jsoniter -ldflags "-s -w" -o dist/app . && \
    mv -v run.sh dist/


# Stage 3
FROM alpine:3.19

RUN apk add --no-cache tini

COPY --from=backend /app/dist/ /app
WORKDIR /app/

VOLUME ["/data"]

EXPOSE 80

ENV PORT=80 \
    DATABASE_TYPE=sqlite \
    DATABASE_DSN=/data/webui/production.db \
    GO_ENV=production \
    GIN_MODE=release \
    APP_BASIC_AUTH_USER= \
    APP_BASIC_AUTH_PASSWORD= \
    APP_BASE_URL= \
    APP_RANCHER_V1_AUTH_URL= \
    APP_PORTAINER_AUTH_URL= \
    APP_SSL_CHALLENGE_BASE_URL= \
    APP_SSL_EMAIL=admin@example.com

ENTRYPOINT ["/sbin/tini", "--", "sh", "run.sh"]

CMD ["start_server"]
