package api

import "github.com/getkin/kin-openapi/openapi3"

type OpenAPISpecLoader func() (*openapi3.T, error)
