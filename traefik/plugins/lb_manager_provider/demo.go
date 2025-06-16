package lb_manager_provider

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/wusuopu/lb_manager_provider/provider"
)

// Config the plugin configuration.
type Config struct {
	PollInterval   string `json:"pollInterval,omitempty"`
	LogFile        string `json:"logFile,omitempty"`
	BaseEndpoint   string `json:"baseEndpoint,omitempty"`		// http://<host>/workspaces/:id/
	QueryParams    string `json:"queryParams,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		PollInterval: "5s", // 5 * time.Second
	}
}

// Provider a simple provider plugin.
type Provider struct {
	name          string
	pollInterval  time.Duration

	config        *Config
	configuration *provider.ConfigurationMarshaler

	cancel func()
}

// New creates a new Provider plugin.
func New(ctx context.Context, config *Config, name string) (*Provider, error) {
	pi, err := time.ParseDuration(config.PollInterval)
	if err != nil {
		return nil, err
	}

	if config.LogFile == "" {
		config.LogFile = os.Getenv("LB_MANAGER_LOG_FILE")
	}
	if config.LogFile == "" {
		config.LogFile = "/tmp/lb_manager_provider.log"
	}
	if config.BaseEndpoint == "" {
		config.BaseEndpoint = os.Getenv("LB_MANAGER_BASE_ENDPOINT")
	}
	config.BaseEndpoint = strings.TrimRight(config.BaseEndpoint, "/")
	if config.QueryParams == "" {
		config.QueryParams = os.Getenv("LB_MANAGER_QUERY_PARAMS")
	}
	config.QueryParams = strings.TrimLeft(config.QueryParams, "?")

	return &Provider{
		name:         name,
		pollInterval: pi,
		config:       config,
		configuration: &provider.ConfigurationMarshaler{
			BaseEndpoint: config.BaseEndpoint,
			QueryParams: config.QueryParams,
		},
	}, nil
}

// Init the provider.
func (p *Provider) Init() error {
	if p.pollInterval <= 0 {
		return fmt.Errorf("poll interval must be greater than 0")
	}

	provider.InitLogger(p.config.LogFile)
	if provider.InitSSLDir() != nil {
		return fmt.Errorf("failed to initialize SSL directory")
	}

	return nil
}

// Provide creates and send dynamic configuration.
func (p *Provider) Provide(cfgChan chan<- json.Marshaler) error {
	ctx, cancel := context.WithCancel(context.Background())
	p.cancel = cancel

	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Print(err)
			}
		}()
		// 插件加载先获取一次配置
		p.configuration.LoadConfiguration()
		cfgChan <- p.configuration

		// 然后每隔几秒获取一次配置
		p.loadConfiguration(ctx, cfgChan)
	}()

	return nil
}

func (p *Provider) loadConfiguration(ctx context.Context, cfgChan chan<- json.Marshaler) {
	ticker := time.NewTicker(p.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			p.configuration.LoadConfiguration()

			cfgChan <- p.configuration
		case <-ctx.Done():
			return
		}
	}
}

// Stop to stop the provider and the related go routines.
func (p *Provider) Stop() error {
	p.cancel()
	provider.CloseLogger()
	return nil
}
