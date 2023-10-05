# Go OIDC Auth proxy

## Overview
It's designed to simulate the authentication layer commonly used for serverless applications.
This proxy allows applications to use OpenID Connect (OIDC) infront of any service. Once authenticated,
the proxy provides the backend services with user claims in the form of X-Claims headers.

This project was born out of experimentation with Azure Container Apps. The goal was to mock
Azure's authentication logic.

This service needs was meant to be run with _ngx_http_auth_request_module_ (example below).

Important: This is not meant for production.

Tested with:

* Keycloak

## Configuration

Set environment variables:

* **GOAP_CLIENT_ID**: OIDC provider client ID
* **GOAP_CLIENT_SECRET**: OIDC provider client secret
* **GOAP_ISSUER: OIDC** issuer
* **GOAP_REDIRECT_URL**: Redirect URL for OIDC
* **GOAP_CUSTOM_TEMPLATE_PATH**: Custom template path.

## Usage
1. Using Docker

```bash
docker run -p 8080:8080 ghcr.io/alesbrelih/oidc-auth-proxy:latest
```

or install CMD

```bash
go install github.com/alesbrelih/oidc-auth-proxy/cmd/oidc-auth-proxy
```

2. NGINX configuration example:


```nginx
# docker embedded dns server
resolver 127.0.0.11 valid=1s;

server {
  location /oidc/ {
    set $api "http://oidc_auth_proxy:8080";
    proxy_pass                               $api;
    proxy_set_header X-Real-IP               $remote_addr;
    proxy_set_header X-Scheme                $scheme;
    proxy_set_header X-Auth-Request-Redirect $request_uri;
  }

  location = /oidc/auth {
    set $api "http://oidc_auth_proxy:8080";
    proxy_pass                        $api;
    proxy_set_header Host             $host;
    proxy_set_header X-Real-IP        $remote_addr;
    proxy_set_header X-Scheme         $scheme;
    # nginx auth_request includes headers but not body
    proxy_set_header Content-Length   "";
    proxy_pass_request_body           off;
  }

  location / {
    auth_request /oidc/auth;
    error_page 401 = /oidc/sign-in;
    
    auth_request_set $claims   $upstream_http_x_claims;
    proxy_set_header X-Claims  $claims;
    
    auth_request_set $auth_cookie $upstream_http_set_cookie;
    add_header Set-Cookie $auth_cookie;

    set $myservice "http://example_service:8080";
    proxy_pass $myservice;
  }
}
```

3. Custom template

By default it uses azure template for claims. To modify this:

* Create new template file `touch /tmp/mytemplate.tmpl`
* Mount file to the docker.
* Set `GOAP_CUSTOM_TEMPLATE_PATH` to mounted location.

Example:

**.env**

`GOAP_CUSTOM_TEMPLATE_PATH=/tmp/template.tmpl`

**docker-compose.yml**

```docker-compose.yml
oidc_auth_proxy:
  container_name: goap_auth
  build:
    target: builder
  env_file:
    - '.env'
    - '.env-secret'
  command: make run/oidc-auth-proxy
  expose:
    - 8080
  volumes:
    - ./:/app:cached
    - ./custom_template_example.tmpl:/tmp/template.tmpl
```

## Testing

1. Run `docker-compose up -d keycloak`. 
2. Reset `example-realm` -> `example-client` authorization secret (it's not persisted when exporting realm and clients).
3. Create a client user that will be used for authentication.
3. Update *GOAP_CLIENT_SECRET* and run `docker compose up`.

First set reset Keycloak client in keycloak admin.

## Contribution

### Prerequisites

```bash
go install github.com/ogen-go/ogen@main
```

