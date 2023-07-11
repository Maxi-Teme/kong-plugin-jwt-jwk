# STAGE-1
# Build the kong plugin
FROM golang:alpine AS plugin-builder

WORKDIR /builder

COPY . .

RUN apk add make
RUN make build

# STAGE-2
# Build kong including the plugin that already build in previous stage
FROM kong:alpine

COPY --from=plugin-builder /builder/jwt-jwk /kong/go-plugins/jwt-jwk

USER kong
ENTRYPOINT ["/docker-entrypoint.sh", "kong", "docker-start"]

EXPOSE 8000
EXPOSE 8001
EXPOSE 8443
EXPOSE 8444

STOPSIGNAL SIGQUIT
HEALTHCHECK --interval=10s --timeout=10s --retries=10 CMD kong health
