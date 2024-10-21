package api

import (
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"go-gateway/cmd/conf"
	"go-gateway/cmd/gwLogger"
	"io"
	"net/http"
	"os"
)

func FetchOpenAPISpec(
	config *conf.Config,
	httpClient *http.Client,
	gwLogger *gwLogger.Logger,
) (*openapi3.T, error) {
	resp, err := httpClient.Get(config.BackendConfig.OpenAPISpecURL)
	if err != nil {
		gwLogger.Error(fmt.Sprintf("Failed to fetch OpenAPI spec: %v", err))
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		gwLogger.Error(fmt.Sprintf("Failed to read OpenAPI spec response: %v", err))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		gwLogger.Error(fmt.Sprintf("Failed to fetch OpenAPI spec: %s", body))
	}

	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromData(body)
	if err != nil {
		gwLogger.Error(fmt.Sprintf("Failed to parse OpenAPI document: %v", err))
		return nil, err
	}

	if err := doc.Validate(loader.Context); err != nil {
		gwLogger.Error(fmt.Sprintf("Invalid OpenAPI document: %v", err))
		return nil, err
	}

	gwLogger.Info(fmt.Sprintf("OpenAPI spec loaded successfully from URL: %s", config.BackendConfig.OpenAPISpecURL))
	return doc, nil
}

func LoadOpenAPISpec(filePath string, gwLogger *gwLogger.Logger) (*openapi3.T, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		gwLogger.Error(fmt.Sprintf("Failed to read OpenAPI spec file: %v", err))
		return nil, err
	}

	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromData(data)
	if err != nil {
		gwLogger.Error(fmt.Sprintf("Failed to parse OpenAPI document: %v", err))
		return nil, err
	}

	if err := doc.Validate(loader.Context); err != nil {
		gwLogger.Error(fmt.Sprintf("Invalid OpenAPI document: %v", err))
		return nil, err
	}

	gwLogger.Info(fmt.Sprintf("OpenAPI spec loaded successfully from file: %s", filePath))
	return doc, nil
}

func GetOpenAPISpecLoader(
	config *conf.Config,
	filePath string,
	httpClient *http.Client,
	gwLogger *gwLogger.Logger,
) OpenAPISpecLoader {
	if config.BackendConfig.OpenAPISpecURL != "" {
		return func() (*openapi3.T, error) {
			return FetchOpenAPISpec(config, httpClient, gwLogger)
		}
	}

	return func() (*openapi3.T, error) {
		return LoadOpenAPISpec(filePath, gwLogger)
	}
}
