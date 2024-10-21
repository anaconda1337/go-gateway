package conf

type AppConf struct {
}

type Conf struct {
	BackendConf Config
	AppConf     AppConf
}

type BackendAPIConf struct {
	URL            string `json:"url" yaml:"url"`
	Port           string `json:"port" yaml:"port"`
	OpenAPISpecURL string `json:"openAPISpecURL" yaml:"openAPISpecURL"`
}

type GatewayConf struct {
	Port           string `json:"port" yaml:"port"`
	TimeoutSeconds int    `json:"timeoutSeconds" yaml:"timeoutSeconds"`
	LogLevel       string `json:"logLevel" yaml:"logLevel"`
	LogFile        string `json:"logFile" yaml:"logFile"`
	LogToFile      bool   `json:"logToFile" yaml:"logToFile"`
}

type Config struct {
	BackendConfig BackendAPIConf `json:"backendAPIConfig" yaml:"backendAPIConfig"`
	GatewayConfig GatewayConf    `json:"gatewayConfig" yaml:"gatewayConfig"`
}
