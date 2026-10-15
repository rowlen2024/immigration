package service

import (
	"encoding/json"
	"errors"
	"log/slog"
	"strings"
	"sync"
	"testing"
	"time"

	"mygo-immigration/backend/internal/dto"
	"mygo-immigration/backend/internal/logging"
	"mygo-immigration/backend/internal/model"
)

const (
	// leadAsyncTestTimeout 限制异步通知测试的等待时间。
	leadAsyncTestTimeout = 5 * time.Second
	// leadNotificationQuietPeriod 用于观察失败提交是否错误启动了通知任务。
	leadNotificationQuietPeriod = 50 * time.Millisecond
)

// leadNotifierFunc 将通知测试函数适配为必需依赖。
type leadNotifierFunc func(lead model.Lead) error

// Notify 将已保存的咨询快照交给测试函数。
func (f leadNotifierFunc) Notify(lead model.Lead) error {
	return f(lead)
}

// notificationLogWriter 通过通道捕获异步日志，避免并发读取缓冲区。
type notificationLogWriter struct {
	// entries 接收已完成序列化的日志内容。
	entries chan []byte
}

// Write 将日志副本发送给测试调用方。
func (w notificationLogWriter) Write(data []byte) (int, error) {
	w.entries <- append([]byte(nil), data...)
	return len(data), nil
}

// TestLeadCreateNotificationIsAsync 验证入库后异步通知不阻塞响应，且消费保存结果的独立快照。
func TestLeadCreateNotificationIsAsync(t *testing.T) {
	release := make(chan struct{})
	unblock := sync.OnceFunc(func() {
		close(release)
	})
	notified := make(chan model.Lead, 1)
	finished := make(chan struct{})
	t.Cleanup(func() {
		unblock()
		select {
		case <-finished:
		case <-time.After(leadAsyncTestTimeout):
			t.Error("通知任务未结束")
		}
	})
	createdAt := time.Now()
	repo := &mockLeadRepo{
		createFn: func(lead *model.Lead) error {
			lead.ID = 42
			lead.CreatedAt = createdAt
			return nil
		},
	}
	svc := NewLeadService(repo, leadNotifierFunc(func(lead model.Lead) error {
		defer close(finished)
		<-release
		notified <- lead
		return nil
	}))
	returned := make(chan *model.Lead, 1)
	createErrors := make(chan error, 1)
	go func() {
		lead, err := svc.Create(&dto.LeadRequest{
			Name:    "<b>张三</b>",
			Phone:   "13800138000",
			Message: "<b>希望了解项目</b>",
		})
		if err != nil {
			createErrors <- err
			return
		}
		returned <- lead
	}()

	select {
	case lead := <-returned:
		lead.Name = "调用方修改后的姓名"
	case err := <-createErrors:
		t.Fatal(err)
	case <-time.After(leadAsyncTestTimeout):
		t.Fatal("提交被邮件通知阻塞")
	}
	// 释放通知器后再检查快照，确保调用方修改不会影响异步读取。
	unblock()
	select {
	case lead := <-notified:
		if lead.ID != 42 || lead.Name != "张三" || lead.Message != "希望了解项目" || !lead.CreatedAt.Equal(createdAt) {
			t.Errorf("通知未消费已保存的咨询快照：%#v", lead)
		}
	case <-time.After(leadAsyncTestTimeout):
		t.Error("未收到咨询快照")
	}
}

// TestLeadCreatePersistenceFailureSkipsNotification 验证入库失败时不触发通知。
func TestLeadCreatePersistenceFailureSkipsNotification(t *testing.T) {
	wantErr := errors.New("database unavailable")
	notified := make(chan struct{}, 1)
	svc := NewLeadService(&mockLeadRepo{
		createFn: func(*model.Lead) error {
			return wantErr
		},
	}, leadNotifierFunc(func(model.Lead) error {
		notified <- struct{}{}
		return nil
	}))
	lead, err := svc.Create(&dto.LeadRequest{
		Name:  "张三",
		Phone: "13800138000",
	})
	if lead != nil || !errors.Is(err, wantErr) {
		t.Fatalf("应返回入库错误：%#v，%v", lead, err)
	}
	select {
	case <-notified:
		t.Fatal("入库失败时不应触发通知")
	case <-time.After(leadNotificationQuietPeriod):
	}
}

// TestLeadCreateNotificationFailureOnlyLogs 验证通知失败仅记录咨询编号和错误，不改变成功结果。
func TestLeadCreateNotificationFailureOnlyLogs(t *testing.T) {
	entries := make(chan []byte, 1)
	previousLogger := logging.Logger
	logging.Logger = slog.New(slog.NewJSONHandler(notificationLogWriter{
		entries: entries,
	}, nil))
	t.Cleanup(func() {
		logging.Logger = previousLogger
	})
	svc := NewLeadService(&mockLeadRepo{
		createFn: func(lead *model.Lead) error {
			lead.ID = 73
			return nil
		},
	}, leadNotifierFunc(func(model.Lead) error {
		return errors.New("SMTP authentication failed")
	}))
	lead, err := svc.Create(&dto.LeadRequest{
		Name:  "张三",
		Phone: "13800138000",
	})
	if err != nil || lead.ID != 73 {
		t.Fatalf("通知失败不应影响成功保存的咨询：%#v，%v", lead, err)
	}
	select {
	case entry := <-entries:
		var fields map[string]any
		if err := json.Unmarshal(entry, &fields); err != nil {
			t.Fatal(err)
		}
		if fields["lead_id"] != float64(73) || fields["level"] != "ERROR" || !strings.Contains(string(entry), "SMTP authentication failed") {
			t.Fatalf("通知失败日志缺少咨询编号或错误：%s", entry)
		}
	case <-time.After(leadAsyncTestTimeout):
		t.Fatal("通知失败未记录日志")
	}
}
