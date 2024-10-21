package main

import (
	"go-gateway/cmd/api"
	"go-gateway/cmd/conf"
	"go-gateway/cmd/gateway"
	logger "go-gateway/cmd/gwLogger"
	"log"
)

func main() {
	// Load gateway config
	config, err := conf.LoadConfig(GatewayConf)
	if err != nil {
		log.Fatalf("Error loading gateway config: %v", err)
	}

	gwLogger, err := logger.NewLogger(config)
	if err != nil {
		panic(err)
	}
	defer gwLogger.Close()

	// Create HTTP client
	httpClient := gateway.HttpClient(config)

	// Load OpenAPI spec
	openAPISpecLoader := api.GetOpenAPISpecLoader(config, OpenApiSpec, httpClient, gwLogger)
	openAPISpec, err := openAPISpecLoader()

	newGateway := gateway.NewGateway(config, httpClient, gwLogger)

	// Setup routes
	if err := newGateway.SetupRoutes(openAPISpec); err != nil {
		log.Fatalf("Error setting up routes: %v", err)
	}

	if err := newGateway.Start(); err != nil {
		log.Fatalf("Error starting newGateway: %v", err)
	}
}
