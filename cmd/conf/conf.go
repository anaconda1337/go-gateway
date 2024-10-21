package conf

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"strings"
)

func LoadConfig(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	config := &Config{
		BackendConfig: BackendAPIConf{},
		GatewayConfig: GatewayConf{},
	}

	switch {
	case strings.HasSuffix(filePath, ".yaml"), strings.HasSuffix(filePath, ".yml"):
		err = yaml.Unmarshal(data, config)
	case strings.HasSuffix(filePath, ".json"):
		err = json.Unmarshal(data, config)
	default:
		return nil, fmt.Errorf("unsupported config file format: %s", filePath)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}
	log.Println("Config loaded successfully from file:", filePath)
	return config, nil
}
