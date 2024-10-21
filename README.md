# GO-GATEWAY

## Description
- Dynamic gateway or reverse proxy with easy configuration for microservices.

## TLDR
- The gateway will load OpenApiSpec either from URL or local file.
- The gateway will parse the methods, endpoints, schemas and create corresponding
- The gateway will start serve as a proxy to your backend.

## Features
- Dynamic routing
- Dynamic schema validation on gateway layer **(WIP)**
- Simple configuration
- Logging

## Installation
1. Clone the repository
```bash
git clone git@github.com:anaconda1337/go-gateway.git
```

2. Edit the yaml/json configuration files
- `gateway.json`
- `gateway.yaml`
- - ^ either should be properly configured
- `openapi.json`
- - ^ this one is needed or set up the `openAPISpecURL` to your API `openapi.json` endpoint so can be fetched.

3. Install the dependencies
```bash
go mod tidy
```
4. Run the Gateway
```bash
cd cmd
go run .
```

## To Do
- [ ] Add tests
- [ ] Add documentation
- [ ] Add middleware support (auth, rate limiting, etc)
- [ ] Add schema validation on request body
- [ ] Add config validation
