FROM golang:1.21.4-bullseye as builder

FROM builder as deployer

WORKDIR /app

COPY ./ .

RUN make build/oidc-auth-proxy

FROM gcr.io/distroless/base AS deployable
USER 65534:65534

WORKDIR /app

COPY --chown=65534:65534 --from=deployer /app/bin/oidc-auth-proxy /app/service

ENTRYPOINT ["/app/service"]
