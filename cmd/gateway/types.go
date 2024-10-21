package gateway

import (
	"github.com/gorilla/mux"
	"go-gateway/cmd/conf"
	"go-gateway/cmd/gwLogger"
	"net/http"
)

type Gateway struct {
	router     *mux.Router
	config     conf.Config
	httpClient *http.Client
	gwLogger   *gwLogger.Logger
}
