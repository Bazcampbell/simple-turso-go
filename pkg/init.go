// pkg/init.go

package simpleturso

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	mu sync.RWMutex
	client = &http.Client{Timeout: 10 * time.Second}
	config *Config
)

func getConfig() (*Config, error) {
	mu.RLock()
	defer mu.RUnlock()
	if config == nil {
		return nil, fmt.Errorf("client not initialised. Call Init() first.")
	}
	return config, nil
}

func Init(cfg *Config) error {	
	if !isValidJWTStructure(cfg.DbKey) {
		return fmt.Errorf("invalid db key")
	}

	if !isValidTursoUrl(cfg.DbUrl) {
		return fmt.Errorf("invalid db url")
	}

	mu.Lock()
	defer mu.Unlock()

	// Convert libsql:// → https:// for HTTP API
	if strings.HasPrefix(cfg.DbUrl, "libsql://") {
		cfg.DbUrl = strings.Replace(cfg.DbUrl, "libsql://", "https://", 1)
	}
	
	config = cfg

	err := test()
	if err != nil {
		config = nil
		return fmt.Errorf("test failed: %w", err) 	
	}

	return nil
}

func test() error {
	req := TursoRequest{
		Statements: []Statement{
			{
				Q:      "SELECT 1",
				Params: []interface{}{},
			},
		},
	}

	_, err := executeTursoRequestWithResponse(req, config)
	return err
}
