package oidc

// Generate schemas:

// Both server and client:
//
//go:generate ogen --clean  --no-client --target api  ../../../go-oidc-proxy.yaml
