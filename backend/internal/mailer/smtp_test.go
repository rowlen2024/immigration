package mailer

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"io"
	"mime"
	"mime/quotedprintable"
	"net"
	"net/http/httptest"
	"net/mail"
	"net/textproto"
	"strings"
	"testing"
	"time"

	"mygo-immigration/backend/internal/config"
)

// smtpTestTimeout 限制本地 SMTP 测试的等待时间。
const smtpTestTimeout = 5 * time.Second

// smtpTestResult 保存本地 SMTP 服务端实际接收的协议数据。
type smtpTestResult struct {
	// auth 是解码后的 PLAIN 认证负载。
	auth string
	// from 是 SMTP 信封发件人命令。
	from string
	// to 是 SMTP 信封收件人命令。
	to string
	// message 是服务端接收的完整 MIME 邮件。
	message []byte
	// err 是本地服务端交互错误。
	err error
}

// newSMTPTestServer 启动可信本地 TLS SMTP 服务，支持模拟认证拒绝。
func newSMTPTestServer(t *testing.T, rejectAuth bool) (*SMTPClient, <-chan smtpTestResult) {
	t.Helper()
	certificateServer := httptest.NewTLSServer(nil)
	certificateServer.Close()
	listener, err := tls.Listen("tcp", "127.0.0.1:0", &tls.Config{
		Certificates: certificateServer.TLS.Certificates,
		MinVersion:   tls.VersionTLS12,
	})
	if err != nil {
		t.Fatalf("创建本地 TLS SMTP 服务失败：%v", err)
	}
	t.Cleanup(func() {
		listener.Close()
	})
	results := make(chan smtpTestResult, 1)
	go func() {
		results <- serveSMTPTestSession(listener, rejectAuth)
	}()

	host, port, err := net.SplitHostPort(listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	client := NewSMTPClient(config.SMTPConfig{
		Host:     host,
		Port:     port,
		Username: "sender@example.com",
		Password: "test-password",
	})
	rootCAs := x509.NewCertPool()
	rootCAs.AddCert(certificateServer.Certificate())
	client.tlsConfig.RootCAs = rootCAs
	return client, results
}

// serveSMTPTestSession 执行一次真实 TLS SMTP 会话并捕获信封和正文。
func serveSMTPTestSession(listener net.Listener, rejectAuth bool) (result smtpTestResult) {
	conn, err := listener.Accept()
	if err != nil {
		result.err = err
		return result
	}
	defer conn.Close()
	if err := conn.SetDeadline(time.Now().Add(smtpTestTimeout)); err != nil {
		result.err = err
		return result
	}
	protocol := textproto.NewConn(conn)
	if err := protocol.PrintfLine("220 localhost ESMTP"); err != nil {
		result.err = err
		return result
	}

	for {
		command, err := protocol.ReadLine()
		if err != nil {
			result.err = err
			return result
		}
		switch {
		case strings.HasPrefix(command, "EHLO "):
			err = protocol.PrintfLine("250-localhost\r\n250 AUTH PLAIN")
		case strings.HasPrefix(command, "AUTH PLAIN "):
			payload, decodeErr := base64.StdEncoding.DecodeString(strings.TrimPrefix(command, "AUTH PLAIN "))
			if decodeErr != nil {
				result.err = decodeErr
				return result
			}
			result.auth = string(payload)
			if rejectAuth {
				err = protocol.PrintfLine("535 authentication rejected")
			} else {
				err = protocol.PrintfLine("235 authentication accepted")
			}
		case strings.HasPrefix(command, "MAIL FROM:"):
			result.from = command
			err = protocol.PrintfLine("250 sender accepted")
		case strings.HasPrefix(command, "RCPT TO:"):
			result.to = command
			err = protocol.PrintfLine("250 recipient accepted")
		case command == "DATA":
			if err = protocol.PrintfLine("354 send message"); err == nil {
				result.message, err = protocol.ReadDotBytes()
			}
			if err == nil {
				err = protocol.PrintfLine("250 message accepted")
			}
		case command == "QUIT":
			result.err = protocol.PrintfLine("221 closing connection")
			return result
		case command == "*":
			// Go SMTP 在认证失败后先取消 AUTH 交互，再退出连接。
			err = protocol.PrintfLine("501 authentication cancelled")
		default:
			err = fmt.Errorf("未预期的 SMTP 命令：%s", command)
		}
		if err != nil {
			result.err = err
			return result
		}
	}
}

// waitSMTPTestResult 等待 SMTP 会话结束，避免遗留测试 goroutine。
func waitSMTPTestResult(t *testing.T, results <-chan smtpTestResult) smtpTestResult {
	t.Helper()
	select {
	case result := <-results:
		return result
	case <-time.After(smtpTestTimeout):
		t.Fatal("等待 SMTP 会话结果超时")
		return smtpTestResult{}
	}
}

// TestSMTPClientSend 验证真实 TLS、认证、信封及中文多行邮件的编码传输。
func TestSMTPClientSend(t *testing.T) {
	client, results := newSMTPTestServer(t, false)
	const recipient = "operations@example.com"
	const subject = "新的联系我们咨询"
	const body = "姓名：张三\n留言：希望了解投资移民项目。\n.开头的行\n" +
		"这是一段超过单行传输长度的中文留言，用于确认 MIME 编码自动折行后不会丢失用户提供的咨询内容。"
	if err := client.Send(recipient, subject, body); err != nil {
		t.Fatalf("发送邮件失败：%v", err)
	}
	result := waitSMTPTestResult(t, results)
	if result.err != nil {
		t.Fatalf("SMTP 服务端失败：%v", result.err)
	}
	if result.auth != "\x00sender@example.com\x00test-password" {
		t.Fatal("SMTP 认证账号或密码未按配置传递")
	}
	if result.from != "MAIL FROM:<sender@example.com>" || result.to != "RCPT TO:<operations@example.com>" {
		t.Fatalf("SMTP 信封不正确：%s，%s", result.from, result.to)
	}
	message, err := mail.ReadMessage(bytes.NewReader(result.message))
	if err != nil {
		t.Fatal(err)
	}
	decodedSubject, err := new(mime.WordDecoder).DecodeHeader(message.Header.Get("Subject"))
	if err != nil || decodedSubject != subject {
		t.Fatalf("邮件主题编码不正确：%q，%v", decodedSubject, err)
	}
	from, err := message.Header.AddressList("From")
	if err != nil || len(from) != 1 || from[0].Address != client.config.Username {
		t.Fatalf("邮件发件人不正确：%v，%v", from, err)
	}
	to, err := message.Header.AddressList("To")
	if err != nil || len(to) != 1 || to[0].Address != recipient {
		t.Fatalf("邮件收件人不正确：%v，%v", to, err)
	}
	if message.Header.Get("Content-Type") != "text/plain; charset=UTF-8" || message.Header.Get("Content-Transfer-Encoding") != "quoted-printable" {
		t.Fatal("邮件未使用 UTF-8 纯文本和 quoted-printable 编码")
	}
	decodedBody, err := io.ReadAll(quotedprintable.NewReader(message.Body))
	if err != nil || strings.ReplaceAll(string(decodedBody), "\r\n", "\n") != body+"\n" {
		t.Fatalf("邮件正文传输不完整：%q，%v", decodedBody, err)
	}
}

// TestSMTPClientAuthenticationFailure 验证认证失败向调用方返回且不发送邮件。
func TestSMTPClientAuthenticationFailure(t *testing.T) {
	client, results := newSMTPTestServer(t, true)
	err := client.Send("operations@example.com", "通知", "正文")
	if err == nil || !strings.Contains(err.Error(), "authenticate SMTP") {
		t.Fatalf("应返回 SMTP 认证失败，实际为：%v", err)
	}
	result := waitSMTPTestResult(t, results)
	if result.err != nil || len(result.message) != 0 {
		t.Fatalf("认证失败后不应发送邮件：%v", result.err)
	}
}

// TestSMTPClientRejectsUntrustedTLS 验证发送器不会跳过 SMTP 服务器证书校验。
func TestSMTPClientRejectsUntrustedTLS(t *testing.T) {
	client, results := newSMTPTestServer(t, false)
	client.tlsConfig.RootCAs = x509.NewCertPool()
	if err := client.Send("operations@example.com", "通知", "正文"); err == nil {
		t.Fatal("不应接受未受信任的 SMTP TLS 证书")
	}
	result := waitSMTPTestResult(t, results)
	if result.auth != "" || len(result.message) != 0 {
		t.Fatal("TLS 校验失败时不应传输凭据或邮件")
	}
}
