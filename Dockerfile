FROM golang:1.21.0-bullseye as builder

FROM builder as deployer

WORKDIR /app

COPY ./code/go/ .

RUN make build/go-oidc-proxy

FROM gcr.io/distroless/base AS deployable
USER 65534:65534

WORKDIR /app

COPY --chown=65534:65534 --from=deployer /app/bin/go-oidc-proxy /app/service

ENTRYPOINT ["/app/service"]
