package main

import (
	"errors"
	"fmt"
	"github.com/ardanlabs/conf"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"time"
)

// WebAPIConfiguration describes the web API configuration. This structure is automatically parsed by
// loadConfiguration and values from flags, environment variable or configuration file will be loaded.
type WebAPIConfiguration struct {
	Config struct {
		Path string `conf:"default:./conf/config.yaml"`
	}
	Web struct {
		APIHost         string        `conf:"default:0.0.0.0:3000,env:CFG_WEB_API_HOST"`
		DebugHost       string        `conf:"default:0.0.0.0:4000,env:CFG_WEB_DEBUG_HOST"`
		ReadTimeout     time.Duration `conf:"default:5s,env:CFG_WEB_READ_TIMEOUT"`
		WriteTimeout    time.Duration `conf:"default:5s,env:CFG_WEB_WRITE_TIMEOUT"`
		ShutdownTimeout time.Duration `conf:"default:5s,env:CFG_WEB_SHUTDOWN_TIMEOUT"`
	}
	Debug bool `conf:"default:false,env:CFG_WEB_DEBUG"`
	DB    struct {
		Filename string `conf:"default:/tmp/decaf.db,env:CFG_WEB_DB_FILENAME"`
	}
	MediaStorage struct {
		UrlPath string `conf:"default:/media/images/*filepath"`
		RootDir string `conf:"default:./media/images,env:CFG_WEB_MEDIA_ROOT_DIR"`
	}
}

// loadConfiguration creates a WebAPIConfiguration starting from flags, environment variables and configuration file.
// It works by loading environment variables first, then update the config using command line flags, finally loading the
// configuration file (specified in WebAPIConfiguration.Config.Path).
// So, CLI parameters will override the environment, and configuration file will override everything.
// Note that the configuration file can be specified only via CLI or environment variable.
func loadConfiguration() (WebAPIConfiguration, error) {
	var cfg WebAPIConfiguration

	// Try to load configuration from environment variables and command line switches
	if err := conf.Parse(os.Args[1:], "CFG", &cfg); err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			usage, err := conf.Usage("CFG", &cfg)
			if err != nil {
				return cfg, fmt.Errorf("generating config usage: %w", err)
			}
			fmt.Println(usage) // nolint:forbidigo
			return cfg, conf.ErrHelpWanted
		}
		return cfg, fmt.Errorf("parsing config: %w", err)
	}

	// Override values from YAML if specified and if it exists (useful in k8s/compose)
	fp, err := os.Open(cfg.Config.Path)
	if err != nil && !os.IsNotExist(err) {
		return cfg, fmt.Errorf("can't read the config file, while it exists: %w", err)
	} else if err == nil {
		yamlFile, err := io.ReadAll(fp)
		if err != nil {
			return cfg, fmt.Errorf("can't read config file: %w", err)
		}
		err = yaml.Unmarshal(yamlFile, &cfg)
		if err != nil {
			return cfg, fmt.Errorf("can't unmarshal config file: %w", err)
		}
		_ = fp.Close()
	}
	return cfg, nil
}
