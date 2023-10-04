# Go OIDC Auth proxy

## Overview
It's designed to simulate the authentication layer commonly used for serverless applications.
This proxy allows applications to use OpenID Connect (OIDC) infront of any service. Once authenticated,
the proxy provides the backend services with user claims in the form of X-Claims headers.

This project was born out of experimentation with Azure Container Apps. The goal was to mock
Azure's authentication logic and provide a flexible and robust authentication solution for serverless applications.


## Configuration

Set environment variables:

* *GOAP_CLIENT_ID*: OIDC provider client ID
* *GOAP_CLIENT_SECRET*: OIDC provider client secret
* *GOAP_ISSUER: OIDC* issuer
* *GOAP_REDIRECT_URL*: Redirect URL for OIDC

## Running
1. Using Docker

```bash
docker run -p 8080:8080 ghcr.io/alesbrelih/go-oidc-auth-proxy:latest
```

2. Install CMD

```bash
go install github.com/alesbrelih/go-oidc-auth-proxy/cmd/go-oidc-auth-proxy
```

### Prerequisites

```bash
go install github.com/ogen-go/ogen@main
```

