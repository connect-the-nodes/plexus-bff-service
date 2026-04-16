package config

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server        ServerConfig        `yaml:"server"`
	Spring        SpringConfig        `yaml:"spring"`
	Management    ManagementConfig    `yaml:"management"`
	Security      SecurityConfig      `yaml:"security"`
	SpringDoc     SpringDocConfig     `yaml:"springdoc"`
	Logging       LoggingConfig       `yaml:"logging"`
	Features      FeaturesConfig      `yaml:"features"`
	AWS           AWSConfig           `yaml:"aws"`
	PropertyLogger PropertyLoggerConfig `yaml:"property-logger"`
	Auth          AuthConfig          `yaml:"auth"`
	Observability ObservabilityConfig `yaml:"observability"`
	activeProfile string
}

type ServerConfig struct {
	Port int `yaml:"port"`
}

func (s ServerConfig) ListenAddress() string {
	port := s.Port
	if port == 0 {
		port = 8080
	}
	return net.JoinHostPort("", strconv.Itoa(port))
}

type SpringConfig struct {
	Application ApplicationConfig `yaml:"application"`
	Profiles    ProfilesConfig    `yaml:"profiles"`
	Security    SpringUserConfig  `yaml:"security"`
	Retry       RetryConfig       `yaml:"retry"`
	Session     SessionConfig     `yaml:"session"`
	Redis       RedisConfig       `yaml:"redis"`
	Data        DataConfig        `yaml:"data"`
}

type ApplicationConfig struct {
	Name string `yaml:"name"`
}

type ProfilesConfig struct {
	Active string `yaml:"active"`
}

type SpringUserConfig struct {
	User UserConfig `yaml:"user"`
}

type UserConfig struct {
	Name     string `yaml:"name"`
	Password string `yaml:"password"`
	Roles    string `yaml:"roles"`
}

type RetryConfig struct {
	MaxAttempts int  `yaml:"maxAttempts"`
	Delay       int  `yaml:"delay"`
	MaxDelay    int  `yaml:"maxDelay"`
	Random      bool `yaml:"random"`
}

type SessionConfig struct {
	StoreType string `yaml:"store-type"`
}

type RedisConfig struct {
	Namespace string `yaml:"namespace"`
}

type DataConfig struct {
	Redis RedisDataConfig `yaml:"redis"`
}

type RedisDataConfig struct {
	Host               string          `yaml:"host"`
	Port               int             `yaml:"port"`
	SSL                RedisSSLConfig  `yaml:"ssl"`
	IAM                RedisIAMConfig  `yaml:"iam"`
	FailFast           bool            `yaml:"fail-fast"`
	UserID             string          `yaml:"userId"`
	ReplicationGroupID string          `yaml:"replicationGroupId"`
	Region             string          `yaml:"region"`
}

type RedisSSLConfig struct {
	Enabled bool `yaml:"enabled"`
}

type RedisIAMConfig struct {
	Enabled bool `yaml:"enabled"`
}

type ManagementConfig struct {
	Health HealthConfig `yaml:"health"`
}

type HealthConfig struct {
	Redis ToggleConfig `yaml:"redis"`
}

type ToggleConfig struct {
	Enabled bool `yaml:"enabled"`
}

type SecurityConfig struct {
	Enabled bool            `yaml:"enabled"`
	JWT     SecurityJWTConfig `yaml:"jwt"`
}

type SecurityJWTConfig struct {
	JWKSetURI      string `yaml:"jwk-set-uri"`
	IssuerURI      string `yaml:"issuer-uri"`
	AuthoritiesClaim string `yaml:"authorities-claim"`
	RolePrefix     string `yaml:"role-prefix"`
}

type SpringDocConfig struct {
	SwaggerUI SwaggerUIConfig `yaml:"swagger-ui"`
}

type SwaggerUIConfig struct {
	Enabled bool   `yaml:"enabled"`
	Path    string `yaml:"path"`
	URL     string `yaml:"url"`
}

type LoggingConfig struct {
	Level map[string]string `yaml:"level"`
}

type FeaturesConfig struct {
	File    string                `yaml:"file"`
	Session FeatureSessionConfig  `yaml:"session"`
}

type FeatureSessionConfig struct {
	Enabled bool `yaml:"enabled"`
}

type AWSConfig struct {
	AppConfig AWSAppConfig `yaml:"app-config"`
}

type AWSAppConfig struct {
	Features AppConfigFeatures `yaml:"features"`
}

type AppConfigFeatures struct {
	Enabled         bool   `yaml:"enabled"`
	ApplicationID   string `yaml:"application-id"`
	EnvironmentID   string `yaml:"environment-id"`
	ConfigurationID string `yaml:"configuration-id"`
}

type PropertyLoggerConfig struct {
	Enabled bool `yaml:"enabled"`
}

type AuthConfig struct {
	Cognito CognitoConfig `yaml:"cognito"`
}

type CognitoConfig struct {
	Enabled              bool   `yaml:"enabled"`
	Domain               string `yaml:"domain"`
	ClientID             string `yaml:"client-id"`
	RedirectURI          string `yaml:"redirect-uri"`
	PostLoginRedirectURI string `yaml:"post-login-redirect-uri"`
	Scopes               string `yaml:"scopes"`
}

type ObservabilityConfig struct {
	Service ObservabilityServiceConfig `yaml:"service"`
}

type ObservabilityServiceConfig struct {
	BaseURL string `yaml:"base-url"`
}

var placeholderPattern = regexp.MustCompile(`\$\{([A-Z0-9_]+)(?::([^}]*))?\}`)

func Load(configDir string) (*Config, error) {
	basePath := filepath.Join(configDir, "application.yml")
	cfg := &Config{}
	if err := loadFile(basePath, cfg); err != nil {
		return nil, err
	}

	profile := os.Getenv("SPRING_PROFILES_ACTIVE")
	if profile == "" {
		profile = cfg.Spring.Profiles.Active
	}
	if profile == "" {
		profile = "local"
	}
	cfg.activeProfile = profile

	for _, current := range strings.Split(profile, ",") {
		current = strings.TrimSpace(current)
		if current == "" {
			continue
		}
		if err := loadFile(filepath.Join(configDir, fmt.Sprintf("application-%s.yml", current)), cfg); err != nil {
			if !os.IsNotExist(err) {
				return nil, err
			}
		}
	}

	if cfg.Features.File == "" {
		cfg.Features.File = "features.yml"
	}
	if cfg.Security.JWT.AuthoritiesClaim == "" {
		cfg.Security.JWT.AuthoritiesClaim = "cognito:groups"
	}
	if cfg.Security.JWT.RolePrefix == "" {
		cfg.Security.JWT.RolePrefix = "ROLE_"
	}
	return cfg, nil
}

func (c *Config) ActiveProfile() string {
	return c.activeProfile
}

func loadFile(path string, target *Config) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	expanded := placeholderPattern.ReplaceAllStringFunc(string(raw), func(match string) string {
		parts := placeholderPattern.FindStringSubmatch(match)
		if len(parts) < 2 {
			return match
		}
		if value := os.Getenv(parts[1]); value != "" {
			return value
		}
		if len(parts) > 2 {
			return parts[2]
		}
		return ""
	})

	return yaml.Unmarshal([]byte(expanded), target)
}
