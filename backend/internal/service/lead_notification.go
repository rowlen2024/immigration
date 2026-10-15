package service

import (
	"fmt"
	"time"

	"mygo-immigration/backend/internal/model"
)

// leadNotificationSubject 是联系咨询通知的固定邮件主题。
const leadNotificationSubject = "新的联系我们咨询"

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

// notificationMailer 提供单收件人纯文本邮件发送能力。
type notificationMailer interface {
	// Send 按原值使用收件人、主题和正文发送邮件。
	Send(to, subject, body string) error
}

// LeadEmailNotifier 读取运营通知邮箱并发送联系咨询邮件。
type LeadEmailNotifier struct {
	// siteConfig 提供运营通知邮箱配置。
	siteConfig siteConfigReader
	// mailer 负责 SMTP 传输。
	mailer notificationMailer
}

// NewLeadEmailNotifier 注入必需的网站设置读取器和邮件发送器。
func NewLeadEmailNotifier(siteConfig siteConfigReader, mailer notificationMailer) *LeadEmailNotifier {
	return &LeadEmailNotifier{
		siteConfig: siteConfig,
		mailer:     mailer,
	}
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

	body := fmt.Sprintf("收到新的联系我们咨询。\n\n咨询编号：%d\n姓名：%s\n电话：%s\n邮箱：%s\n意向项目：%s\n提交时间：%s\n\n留言：\n%s\n",
		lead.ID,
		lead.Name,
		lead.Phone,
		lead.Email,
		lead.InterestedProject,
		lead.CreatedAt.Format(time.RFC3339),
		lead.Message,
	)
	return n.mailer.Send(cfg.NotificationEmail, leadNotificationSubject, body)
}
