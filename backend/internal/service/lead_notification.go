package service

import (
	"bytes"
	_ "embed"
	"fmt"
	"html"
	"html/template"
	"strings"
	"time"

	"mygo-immigration/backend/internal/model"
)

// leadNotificationSubject 是联系咨询通知的固定邮件主题。
const leadNotificationSubject = "北极星移民官网新咨询"

// leadNotificationTemplateSource 随后端二进制嵌入咨询通知邮件模板。
//
//go:embed templates/lead_notification.html
var leadNotificationTemplateSource string

// leadNotificationTemplate 在初始化时解析固定模板，语法错误直接暴露。
var leadNotificationTemplate = template.Must(template.New("lead-notification").Parse(leadNotificationTemplateSource))

// LeadNotifier 定义咨询保存成功后的通知能力。
type LeadNotifier interface {
	// Notify 处理已保存咨询的通知，失败由调用方记录日志。
	Notify(lead model.Lead) error
}

// siteConfigReader 提供通知器需要的网站设置读取能力。
type siteConfigReader interface {
	// GetSiteConfig 返回当前网站设置。
	GetSiteConfig() (*SiteConfig, error)
}

// notificationProjectReader 提供不加载关联数据的项目名称查询能力。
type notificationProjectReader interface {
	// FindBySlugsLight 按原始 slug 查询项目记录。
	FindBySlugsLight(slugs []string) ([]model.Project, error)
}

// notificationMailer 提供单收件人 HTML 邮件发送能力。
type notificationMailer interface {
	// SendHTML 按原值使用收件人、主题和 HTML 正文发送邮件。
	SendHTML(to, subject, body string) error
}

// LeadEmailNotifier 读取运营通知邮箱并发送联系咨询邮件。
type LeadEmailNotifier struct {
	// siteConfig 提供运营通知邮箱配置。
	siteConfig siteConfigReader
	// projects 提供意向项目对应的名称。
	projects notificationProjectReader
	// mailer 负责 SMTP 传输。
	mailer notificationMailer
}

// NewLeadEmailNotifier 注入必需的网站设置读取器、项目读取器和邮件发送器。
func NewLeadEmailNotifier(siteConfig siteConfigReader, projects notificationProjectReader, mailer notificationMailer) *LeadEmailNotifier {
	return &LeadEmailNotifier{
		siteConfig: siteConfig,
		projects:   projects,
		mailer:     mailer,
	}
}

// leadNotificationView 保存邮件模板所需的咨询快照和展示字段。
type leadNotificationView struct {
	// Lead 是已保存且补全项目名称的咨询快照。
	model.Lead
	// SubmittedAt 使用保存结果的时区展示 datetime 格式。
	SubmittedAt string
	// MessageLines 将留言按换行拆分，交由模板逐行转义。
	MessageLines []string
}

// Notify 向已配置的运营邮箱发送咨询快照，未配置时直接结束。
func (n *LeadEmailNotifier) Notify(lead model.Lead) error {
	cfg, err := n.siteConfig.GetSiteConfig()
	if err != nil {
		return fmt.Errorf("read lead notification config: %w", err)
	}
	if cfg.NotificationEmail == "" {
		return nil
	}

	if lead.InterestedProject != "" {
		projects, err := n.projects.FindBySlugsLight([]string{lead.InterestedProject})
		if err != nil {
			return fmt.Errorf("read lead notification project: %w", err)
		}
		if len(projects) > 0 {
			lead.ProjectName = projects[0].Name
		}
	}

	// 保存入口的净化结果包含 HTML 实体，先还原一层，再由模板统一安全转义。
	lead.Name = html.UnescapeString(lead.Name)
	lead.Message = html.UnescapeString(lead.Message)
	view := leadNotificationView{
		Lead:         lead,
		SubmittedAt:  lead.CreatedAt.Format(time.DateTime),
		MessageLines: strings.Split(lead.Message, "\n"),
	}
	var body bytes.Buffer
	if err := leadNotificationTemplate.Execute(&body, view); err != nil {
		return fmt.Errorf("render lead notification: %w", err)
	}
	return n.mailer.SendHTML(cfg.NotificationEmail, leadNotificationSubject, body.String())
}
