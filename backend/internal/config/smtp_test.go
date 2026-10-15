package config

import "testing"

// TestLoadSMTP 验证 SMTP 参数从环境变量加载且默认使用腾讯企业邮箱 TLS 端点。
func TestLoadSMTP(t *testing.T) {
	t.Setenv("SMTP_HOST", "")
	t.Setenv("SMTP_PORT", "")
	t.Setenv("SMTP_USERNAME", "")
	t.Setenv("SMTP_PASSWORD", "")

	cfg := Load().SMTP
	if cfg.Host != DefaultSMTPHost || cfg.Port != DefaultSMTPPort {
		t.Fatalf("SMTP 默认端点不正确：%s:%s", cfg.Host, cfg.Port)
	}
	if cfg.Username != "" || cfg.Password != "" {
		t.Fatal("SMTP 账号和密码不应包含默认凭据")
	}

	t.Setenv("SMTP_HOST", "smtp.example.com")
	t.Setenv("SMTP_PORT", "2465")
	t.Setenv("SMTP_USERNAME", "sender@example.com")
	t.Setenv("SMTP_PASSWORD", "test-password#suffix")
	cfg = Load().SMTP
	if cfg.Host != "smtp.example.com" || cfg.Port != "2465" || cfg.Username != "sender@example.com" || cfg.Password != "test-password#suffix" {
		t.Fatal("SMTP 环境变量未完整加载")
	}
}
