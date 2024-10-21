package gateway

import (
	"bytes"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gorilla/mux"
	"go-gateway/cmd/conf"
	"go-gateway/cmd/gwLogger"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func HttpClient(config *conf.Config) *http.Client {
	return &http.Client{
		Timeout: time.Duration(config.GatewayConfig.TimeoutSeconds) * time.Second,
	}
}

func NewGateway(config *conf.Config, httpClient *http.Client, gwLogger *gwLogger.Logger) *Gateway {
	return &Gateway{
		router:     mux.NewRouter(),
		config:     *config,
		httpClient: httpClient,
		gwLogger:   gwLogger,
	}
}

func getOperations(pathItem *openapi3.PathItem) map[string]*openapi3.Operation {
	methodMap := map[string]*openapi3.Operation{
		"GET":     pathItem.Get,
		"POST":    pathItem.Post,
		"PUT":     pathItem.Put,
		"PATCH":   pathItem.Patch,
		"DELETE":  pathItem.Delete,
		"OPTIONS": pathItem.Options,
		"HEAD":    pathItem.Head,
		"TRACE":   pathItem.Trace,
	}

	operations := make(map[string]*openapi3.Operation)
	for method, op := range methodMap {
		if op != nil {
			operations[method] = op
		}
	}
	return operations
}

func (g *Gateway) SetupRoutes(openAPIDoc *openapi3.T) error {
	for path, pathItem := range openAPIDoc.Paths.Map() {
		routePath := convertPath(path)

		for method, _ := range getOperations(pathItem) {
			handler := g.createHandler(method)
			g.router.HandleFunc(routePath, handler).Methods(method)
			g.gwLogger.Info(fmt.Sprintf("Added route: %s %s", method, routePath))
		}
	}
	return nil
}

func convertPath(path string) string {
	segments := strings.Split(path, "/")
	for i, segment := range segments {
		if strings.HasPrefix(segment, "{") && strings.HasSuffix(segment, "}") {
			param := segment[1 : len(segment)-1]
			segments[i] = fmt.Sprintf("{%s:[^/]+}", param)
		}
	}
	return strings.Join(segments, "/")
}

func (g *Gateway) createHandler(method string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Incoming request: %s %s", r.Method, r.URL.Path)

		if r.Method != method {
			g.gwLogger.Error(fmt.Sprintf("Method Not Allowed: %s", r.Method))
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		backendURL := fmt.Sprintf("%s:%s%s",
			g.config.BackendConfig.URL,
			g.config.BackendConfig.Port,
			r.URL.Path,
		)
		if r.URL.RawQuery != "" {
			backendURL += "?" + r.URL.RawQuery
		}

		g.gwLogger.Info(fmt.Sprintf("Forwarding request to backend: %s", backendURL))

		body, err := io.ReadAll(r.Body)
		if err != nil {
			g.gwLogger.Error(fmt.Sprintf("Error reading request body: %v", err))
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}
		err = r.Body.Close()
		if err != nil {
			g.gwLogger.Error(fmt.Sprintf("Error closing request body: %v", err))
			return
		}

		backendReq, err := http.NewRequest(r.Method, backendURL, bytes.NewBuffer(body))
		if err != nil {
			g.gwLogger.Error(fmt.Sprintf("Error creating backend request: %v", err))
			http.Error(w, "Error creating backend request", http.StatusInternalServerError)
			return
		}

		for key, values := range r.Header {
			for _, value := range values {
				backendReq.Header.Add(key, value)
			}
		}

		log.Printf("Sending request to backend...")
		resp, err := g.httpClient.Do(backendReq)
		if err != nil {
			g.gwLogger.Error(fmt.Sprintf("Error forwarding request to backend: %v", err))
			if urlErr, ok := err.(*url.Error); ok {
				if urlErr.Timeout() {
					http.Error(w, "Backend request timed out", http.StatusGatewayTimeout)
					return
				}
				if urlErr.Temporary() {
					http.Error(w, "Temporary backend error", http.StatusServiceUnavailable)
					return
				}
			}
			http.Error(w, "Error forwarding request to backend", http.StatusBadGateway)
			return
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				g.gwLogger.Error(fmt.Sprintf("Error closing response body: %v", err))
			}
		}(resp.Body)

		g.gwLogger.Info(fmt.Sprintf("Received response from backend: %s", resp.Status))

		for key, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}

		w.WriteHeader(resp.StatusCode)

		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			g.gwLogger.Error(fmt.Sprintf("Error reading response body: %v", err))
			return
		}

		if _, err := w.Write(responseBody); err != nil {
			g.gwLogger.Error(fmt.Sprintf("Error writing response body: %v", err))
		}
	}
}

func (g *Gateway) Start() error {
	g.gwLogger.Info(fmt.Sprintf("Starting gateway on port %s", g.config.GatewayConfig.Port))
	return http.ListenAndServe(":"+g.config.GatewayConfig.Port, g.router)
}
