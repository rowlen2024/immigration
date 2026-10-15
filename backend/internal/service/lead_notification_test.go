package service

import (
	"encoding/json"
	"errors"
	"html"
	"strings"
	"testing"
	"time"

	"mygo-immigration/backend/internal/model"
)

// notificationMailerFunc 将邮件发送测试逻辑适配为通知器依赖。
type notificationMailerFunc func(to, subject, body string) error

// SendHTML 将 HTML 邮件参数原样交给测试逻辑。
func (f notificationMailerFunc) SendHTML(to, subject, body string) error {
	return f(to, subject, body)
}

// notificationProjectReaderFunc 将轻量项目查询测试逻辑适配为必需依赖。
type notificationProjectReaderFunc func(slugs []string) ([]model.Project, error)

// FindBySlugsLight 将原始 slug 交给测试逻辑。
func (f notificationProjectReaderFunc) FindBySlugsLight(slugs []string) ([]model.Project, error) {
	return f(slugs)
}

// unexpectedNotificationProjects 为不应查询项目的场景提供必需依赖。
func unexpectedNotificationProjects(t *testing.T) notificationProjectReaderFunc {
	t.Helper()
	return func([]string) ([]model.Project, error) {
		t.Error("当前通知场景不应查询项目")
		return nil, nil
	}
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
		CreatedAt:         time.Date(2026, time.October, 15, 10, 30, 0, 0, time.FixedZone("UTC+08:00", int((8*time.Hour)/time.Second))),
	}
	const projectName = "希腊黄金签证"
	projectCalls := 0
	projects := notificationProjectReaderFunc(func(slugs []string) ([]model.Project, error) {
		projectCalls++
		if len(slugs) != 1 || slugs[0] != lead.InterestedProject {
			t.Fatalf("项目查询未按原值使用意向项目 slug：%v", slugs)
		}
		return []model.Project{
			{
				Slug: lead.InterestedProject,
				Name: projectName,
			},
		}, nil
	})
	sendCalls := 0
	notifier := NewLeadEmailNotifier(settings, projects, notificationMailerFunc(func(to, subject, body string) error {
		sendCalls++
		if to != "Ops.Team+leads@example.com" || subject != "北极星移民官网新咨询" {
			t.Fatalf("通知收件人或主题不正确：%q，%q", to, subject)
		}
		visibleBody := html.UnescapeString(body)
		for _, value := range []string{"<!doctype html>", "#42", lead.Name, lead.Phone, lead.Email, projectName, "希望了解项目。<br>请电话联系。", "2026-10-15 10:30:00"} {
			if !strings.Contains(visibleBody, value) {
				t.Errorf("邮件正文缺少咨询事实：%q", value)
			}
		}
		if strings.Contains(body, lead.InterestedProject) || strings.Contains(body, lead.CreatedAt.Format(time.RFC3339)) {
			t.Error("邮件不应展示项目 slug 或 RFC3339 提交时间")
		}
		return nil
	}))
	if err := notifier.Notify(lead); err != nil {
		t.Fatal(err)
	}
	if sendCalls != 1 {
		t.Fatalf("应发送一次通知，实际为 %d 次", sendCalls)
	}
	if projectCalls != 1 {
		t.Fatalf("应轻量查询一次项目，实际为 %d 次", projectCalls)
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
			notifier := NewLeadEmailNotifier(settings, unexpectedNotificationProjects(t), notificationMailerFunc(func(string, string, string) error {
				t.Error("未配置通知邮箱时不应调用发送器")
				return nil
			}))
			if err := notifier.Notify(model.Lead{
				InterestedProject: "project-a",
			}); err != nil {
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
	notifier := NewLeadEmailNotifier(settings, unexpectedNotificationProjects(t), notificationMailerFunc(func(string, string, string) error {
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
	notifier := NewLeadEmailNotifier(settings, unexpectedNotificationProjects(t), notificationMailerFunc(func(string, string, string) error {
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
	notifier := NewLeadEmailNotifier(settings, unexpectedNotificationProjects(t), notificationMailerFunc(func(to, subject, body string) error {
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

// TestLeadEmailNotifierOptionalProject 验证未选项目和已删除项目的展示不会回退为 slug。
func TestLeadEmailNotifierOptionalProject(t *testing.T) {
	tests := []struct {
		// name 标识项目场景。
		name string
		// slug 是保存的意向项目原值。
		slug string
		// wantProject 是邮件中的项目展示文本。
		wantProject string
		// wantCalls 是预期的轻量查询次数。
		wantCalls int
	}{
		{
			name:        "未选项目",
			wantProject: "未指定",
			wantCalls:   0,
		},
		{
			name:        "项目已删除",
			slug:        "deleted-project",
			wantProject: "未找到对应项目",
			wantCalls:   1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			settings := &HomeConfigService{
				repo: &mockHomeConfigRepo{
					configs: map[string]json.RawMessage{
						"site": json.RawMessage(`{"notification_email":"operations@example.com"}`),
					},
				},
			}
			projectCalls := 0
			projects := notificationProjectReaderFunc(func(slugs []string) ([]model.Project, error) {
				projectCalls++
				if len(slugs) != 1 || slugs[0] != tt.slug {
					t.Errorf("项目查询参数不正确：%v", slugs)
				}
				return nil, nil
			})
			sendCalls := 0
			notifier := NewLeadEmailNotifier(settings, projects, notificationMailerFunc(func(to, subject, body string) error {
				sendCalls++
				if !strings.Contains(body, tt.wantProject) || (tt.slug != "" && strings.Contains(body, tt.slug)) {
					t.Errorf("意向项目展示不正确：%s", body)
				}
				return nil
			}))
			if err := notifier.Notify(model.Lead{
				InterestedProject: tt.slug,
			}); err != nil {
				t.Fatal(err)
			}
			if projectCalls != tt.wantCalls || sendCalls != 1 {
				t.Fatalf("查询或发送次数不正确：查询 %d 次，发送 %d 次", projectCalls, sendCalls)
			}
		})
	}
}

// TestLeadEmailNotifierProjectFailure 验证项目查询失败交给异步任务记录且不发送不完整邮件。
func TestLeadEmailNotifierProjectFailure(t *testing.T) {
	settings := &HomeConfigService{
		repo: &mockHomeConfigRepo{
			configs: map[string]json.RawMessage{
				"site": json.RawMessage(`{"notification_email":"operations@example.com"}`),
			},
		},
	}
	wantErr := errors.New("project query failed")
	projects := notificationProjectReaderFunc(func([]string) ([]model.Project, error) {
		return nil, wantErr
	})
	notifier := NewLeadEmailNotifier(settings, projects, notificationMailerFunc(func(string, string, string) error {
		t.Error("项目查询失败时不应发送邮件")
		return nil
	}))
	if err := notifier.Notify(model.Lead{
		InterestedProject: "project-a",
	}); !errors.Is(err, wantErr) {
		t.Fatalf("应返回项目查询错误，实际为：%v", err)
	}
}

// TestLeadEmailNotifierHTMLText 验证动态文本安全转义、保存实体不重复转义且多行留言保留空行。
func TestLeadEmailNotifierHTMLText(t *testing.T) {
	settings := &HomeConfigService{
		repo: &mockHomeConfigRepo{
			configs: map[string]json.RawMessage{
				"site": json.RawMessage(`{"notification_email":"operations@example.com"}`),
			},
		},
	}
	projects := notificationProjectReaderFunc(func([]string) ([]model.Project, error) {
		return []model.Project{
			{
				Slug: "project-a",
				Name: "A&B <项目>",
			},
		}, nil
	})
	sendCalls := 0
	notifier := NewLeadEmailNotifier(settings, projects, notificationMailerFunc(func(to, subject, body string) error {
		sendCalls++
		for _, value := range []string{"王 &amp; 李", "A&amp;B &lt;项目&gt;", "第一行 &amp; 联系<br><br>&lt;img src=x onerror=alert(1)&gt;"} {
			if !strings.Contains(body, value) {
				t.Errorf("HTML 邮件未正确展示文本：%q", value)
			}
		}
		if strings.Contains(body, "<img") || strings.Contains(body, "<项目>") || strings.Contains(body, "&amp;amp;") {
			t.Error("用户文本被当作 HTML 或发生重复实体转义")
		}
		return nil
	}))
	if err := notifier.Notify(model.Lead{
		Name:              leadSanitizer.Sanitize("王 & 李"),
		InterestedProject: "project-a",
		Message:           leadSanitizer.Sanitize("<b>第一行 & 联系</b>\n\n&lt;img src=x onerror=alert(1)&gt;"),
	}); err != nil {
		t.Fatal(err)
	}
	if sendCalls != 1 {
		t.Fatalf("应发送一次 HTML 邮件，实际为 %d 次", sendCalls)
	}
}
