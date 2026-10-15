package config

import (
	"bufio"
	"os"
	"strings"
	"time"
)

const (
	// DefaultSMTPHost 是腾讯企业邮箱的 SMTP 服务地址。
	DefaultSMTPHost = "smtp.exmail.qq.com"
	// DefaultSMTPPort 是腾讯企业邮箱的隐式 TLS 端口。
	DefaultSMTPPort = "465"
)

// SMTPConfig 保存通过环境变量配置的发件账号和 TLS SMTP 连接参数。
type SMTPConfig struct {
	// Host 是 SMTP 服务器域名。
	Host string
	// Port 是 SMTP 服务器的隐式 TLS 端口。
	Port string
	// Username 同时作为 SMTP 认证账号和发件邮箱。
	Username string
	// Password 是邮箱密码或 SMTP 授权码。
	Password string
}

func init() {
	loadEnvFile(".env")
}

func loadEnvFile(path string) {
	f, err := os.Open(path)
	if err != nil {
		// try parent directory
		f, err = os.Open("../" + path)
		if err != nil {
			return
		}
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		if os.Getenv(key) == "" {
			os.Setenv(key, val)
		}
	}
}

type Config struct {
	ServerPort       string
	DBHost           string
	DBPort           string
	DBUser           string
	DBPassword       string
	DBName           string
	JWTSecret        string
	JWTAccessExpiry  time.Duration
	JWTRefreshExpiry time.Duration
	CORSOrigin       string
	Env              string
	LogLevel         string
	// SMTP 是联系咨询通知使用的服务端邮件配置。
	SMTP SMTPConfig
}

func Load() *Config {
	return &Config{
		ServerPort:       getEnv("SERVER_PORT", "8080"),
		DBHost:           getEnv("DB_HOST", "localhost"),
		DBPort:           getEnv("DB_PORT", "3306"),
		DBUser:           getEnv("DB_USER", "mygo"),
		DBPassword:       getEnv("DB_PASSWORD", "mygo123"),
		DBName:           getEnv("DB_NAME", "mygo_immigration"),
		JWTSecret:        getEnv("JWT_SECRET", "change-me-in-production"),
		JWTAccessExpiry:  getDurationEnv("JWT_ACCESS_EXPIRY", 15*time.Minute),
		JWTRefreshExpiry: getDurationEnv("JWT_REFRESH_EXPIRY", 168*time.Hour),
		CORSOrigin:       getEnv("CORS_ORIGIN", "http://localhost:3000"),
		Env:              getEnv("ENV", "development"),
		LogLevel:         getEnv("LOG_LEVEL", "info"),
		SMTP: SMTPConfig{
			Host:     getEnv("SMTP_HOST", DefaultSMTPHost),
			Port:     getEnv("SMTP_PORT", DefaultSMTPPort),
			Username: getEnv("SMTP_USERNAME", ""),
			Password: getEnv("SMTP_PASSWORD", ""),
		},
	}
}

func (c *Config) DSN() string {
	return c.DBUser + ":" + c.DBPassword + "@tcp(" + c.DBHost + ":" + c.DBPort + ")/" + c.DBName + "?charset=utf8mb4&parseTime=True&loc=Local&tls=skip-verify"
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getDurationEnv(key string, defaultVal time.Duration) time.Duration {
	if val := os.Getenv(key); val != "" {
		if d, err := time.ParseDuration(val); err == nil {
			return d
		}
	}
	return defaultVal
}
