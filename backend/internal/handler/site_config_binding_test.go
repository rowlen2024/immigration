package handler

import (
	"net/http/httptest"
	"testing"

	"mygo-immigration/backend/internal/service"

	"github.com/gin-gonic/gin"
)

// TestSiteConfigNotificationEmailBinding 验证设置入口接受空值和合法邮箱，拒绝非法原始输入。
func TestSiteConfigNotificationEmailBinding(t *testing.T) {
	tests := []struct {
		// email 是提交给网站设置入口的原始邮箱值。
		email string
		// valid 表示该值是否符合设置接口契约。
		valid bool
	}{
		{
			email: "",
			valid: true,
		},
		{
			email: "Ops.Team+leads@example.com",
			valid: true,
		},
		{
			email: "not-an-email",
			valid: false,
		},
		{
			email: " operations@example.com",
			valid: false,
		},
		{
			email: "operations@example.com ",
			valid: false,
		},
		{
			email: "operations@example.com\r\nBcc: other@example.com",
			valid: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.email, func(t *testing.T) {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Request = makePostRequest("/api/v1/admin/site-config", map[string]string{
				"notification_email": tt.email,
			})
			var cfg service.SiteConfig
			err := c.ShouldBindJSON(&cfg)
			if (err == nil) != tt.valid {
				t.Fatalf("邮箱绑定结果不符合契约：%q，%v", tt.email, err)
			}
			if cfg.NotificationEmail != tt.email {
				t.Fatalf("入口不应规范化原始邮箱：%q", cfg.NotificationEmail)
			}
		})
	}
}
