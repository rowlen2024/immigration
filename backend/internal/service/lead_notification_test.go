package service

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"mygo-immigration/backend/internal/model"
)

// notificationMailerFunc 将邮件发送测试逻辑适配为通知器依赖。
type notificationMailerFunc func(to, subject, body string) error

// Send 将协议参数原样交给测试逻辑。
func (f notificationMailerFunc) Send(to, subject, body string) error {
	return f(to, subject, body)
}

// TestLeadEmailNotifierConfigured 验证通知邮箱与客服邮箱独立且邮件包含保存后的咨询事实。
func TestLeadEmailNotifierConfigured(t *testing.T) {
	settings := &HomeConfigService{
		repo: &mockHomeConfigRepo{
			configs: map[string]json.RawMessage{
				"site": json.RawMessage(`{"contact_email":"support@example.com","notification_email":"Ops.Team+leads@example.com"}`),
			},
		},
	}
	lead := model.Lead{
		ID:                42,
		Name:              "张三",
		Phone:             "+86 13800138000",
		Email:             "customer@example.com",
		InterestedProject: "project-a",
		Message:           "希望了解项目。\n请电话联系。",
		CreatedAt:         time.Date(2026, time.October, 15, 10, 30, 0, 0, time.UTC),
	}
	sendCalls := 0
	notifier := NewLeadEmailNotifier(settings, notificationMailerFunc(func(to, subject, body string) error {
		sendCalls++
		if to != "Ops.Team+leads@example.com" || subject != leadNotificationSubject {
			t.Fatalf("通知收件人或主题不正确：%q，%q", to, subject)
		}
		for _, value := range []string{"咨询编号：42", lead.Name, lead.Phone, lead.Email, lead.InterestedProject, lead.Message, lead.CreatedAt.Format(time.RFC3339)} {
			if !strings.Contains(body, value) {
				t.Errorf("邮件正文缺少咨询事实：%q", value)
			}
		}
		return nil
	}))
	if err := notifier.Notify(lead); err != nil {
		t.Fatal(err)
	}
	if sendCalls != 1 {
		t.Fatalf("应发送一次通知，实际为 %d 次", sendCalls)
	}
}

// TestLeadEmailNotifierWithoutRecipient 验证未配置、旧配置和空通知邮箱均不发送邮件。
func TestLeadEmailNotifierWithoutRecipient(t *testing.T) {
	tests := []struct {
		// name 标识配置场景。
		name string
		// configs 保存该场景的网站配置。
		configs map[string]json.RawMessage
	}{
		{
			name:    "没有网站配置",
			configs: map[string]json.RawMessage{},
		},
		{
			name: "旧配置仅有客服邮箱",
			configs: map[string]json.RawMessage{
				"site": json.RawMessage(`{"contact_email":"support@example.com"}`),
			},
		},
		{
			name: "通知邮箱为空",
			configs: map[string]json.RawMessage{
				"site": json.RawMessage(`{"notification_email":""}`),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			settings := &HomeConfigService{
				repo: &mockHomeConfigRepo{
					configs: tt.configs,
				},
			}
			notifier := NewLeadEmailNotifier(settings, notificationMailerFunc(func(string, string, string) error {
				t.Error("未配置通知邮箱时不应调用发送器")
				return nil
			}))
			if err := notifier.Notify(model.Lead{}); err != nil {
				t.Fatal(err)
			}
		})
	}
}

// TestLeadEmailNotifierSendFailure 验证发送错误返回给异步任务且不重试。
func TestLeadEmailNotifierSendFailure(t *testing.T) {
	settings := &HomeConfigService{
		repo: &mockHomeConfigRepo{
			configs: map[string]json.RawMessage{
				"site": json.RawMessage(`{"notification_email":"operations@example.com"}`),
			},
		},
	}
	wantErr := errors.New("SMTP rejected")
	sendCalls := 0
	notifier := NewLeadEmailNotifier(settings, notificationMailerFunc(func(string, string, string) error {
		sendCalls++
		return wantErr
	}))
	if err := notifier.Notify(model.Lead{}); !errors.Is(err, wantErr) {
		t.Fatalf("发送失败应返回原始错误，实际为：%v", err)
	}
	if sendCalls != 1 {
		t.Fatalf("发送失败不应重试，实际调用 %d 次", sendCalls)
	}
}

// TestLeadEmailNotifierConfigFailure 验证无法解析设置时返回错误且不调用邮件发送器。
func TestLeadEmailNotifierConfigFailure(t *testing.T) {
	settings := &HomeConfigService{
		repo: &mockHomeConfigRepo{
			configs: map[string]json.RawMessage{
				"site": json.RawMessage(`{`),
			},
		},
	}
	notifier := NewLeadEmailNotifier(settings, notificationMailerFunc(func(string, string, string) error {
		t.Error("配置解析失败时不应发送邮件")
		return nil
	}))
	if err := notifier.Notify(model.Lead{}); err == nil {
		t.Fatal("应返回网站配置解析错误")
	}
}

// TestSiteConfigNotificationEmailRoundTrip 验证通知邮箱独立保存、更新和清空，并在后续通知中生效。
func TestSiteConfigNotificationEmailRoundTrip(t *testing.T) {
	settings := &HomeConfigService{
		repo: &mockHomeConfigRepo{
			configs: map[string]json.RawMessage{},
		},
	}
	var recipients []string
	notifier := NewLeadEmailNotifier(settings, notificationMailerFunc(func(to, subject, body string) error {
		recipients = append(recipients, to)
		return nil
	}))
	for _, recipient := range []string{"operations@example.com", "new-operations@example.com", ""} {
		if err := settings.UpdateSiteConfig(&SiteConfig{
			ContactEmail:      "support@example.com",
			NotificationEmail: recipient,
		}); err != nil {
			t.Fatal(err)
		}
		cfg, err := settings.GetSiteConfig()
		if err != nil || cfg.NotificationEmail != recipient || cfg.ContactEmail != "support@example.com" {
			t.Fatalf("通知邮箱保存后读取不一致：%#v，%v", cfg, err)
		}
		if err := notifier.Notify(model.Lead{}); err != nil {
			t.Fatal(err)
		}
	}
	if len(recipients) != 2 || recipients[0] != "operations@example.com" || recipients[1] != "new-operations@example.com" {
		t.Fatalf("通知收件人未随设置更新生效：%v", recipients)
	}
}
