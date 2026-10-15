// Package mailer 提供服务端 TLS SMTP 邮件发送能力。
package mailer

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"mime"
	"mime/quotedprintable"
	"net"
	"net/mail"
	"net/smtp"
	"time"

	"mygo-immigration/backend/internal/config"
)

// smtpTimeout 限制建立连接以及后续 SMTP 交互的等待时间。
const smtpTimeout = 10 * time.Second

// SMTPClient 使用同一邮箱账号认证并发送邮件。
type SMTPClient struct {
	// config 保存发件账号和 SMTP 连接参数。
	config config.SMTPConfig
	// tlsConfig 定义服务器证书校验和最低 TLS 版本。
	tlsConfig *tls.Config
}

// NewSMTPClient 构造使用隐式 TLS、校验证书的 SMTP 客户端。
func NewSMTPClient(cfg config.SMTPConfig) *SMTPClient {
	return &SMTPClient{
		config: cfg,
		tlsConfig: &tls.Config{
			ServerName: cfg.Host,
			MinVersion: tls.VersionTLS12,
		},
	}
}

// Send 向单个收件人发送 UTF-8 纯文本邮件，协议或网络失败直接返回。
func (s *SMTPClient) Send(to, subject, body string) error {
	dialer := &net.Dialer{
		Timeout: smtpTimeout,
	}
	conn, err := tls.DialWithDialer(dialer, "tcp", net.JoinHostPort(s.config.Host, s.config.Port), s.tlsConfig)
	if err != nil {
		return fmt.Errorf("connect SMTP: %w", err)
	}
	defer conn.Close()
	if err := conn.SetDeadline(time.Now().Add(smtpTimeout)); err != nil {
		return fmt.Errorf("set SMTP deadline: %w", err)
	}

	client, err := smtp.NewClient(conn, s.config.Host)
	if err != nil {
		return fmt.Errorf("create SMTP client: %w", err)
	}
	defer client.Close()
	if err := client.Auth(smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)); err != nil {
		return fmt.Errorf("authenticate SMTP: %w", err)
	}
	if err := client.Mail(s.config.Username); err != nil {
		return fmt.Errorf("set SMTP sender: %w", err)
	}
	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("set SMTP recipient: %w", err)
	}

	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("start SMTP message: %w", err)
	}
	if _, err := writer.Write(encodeMessage(s.config.Username, to, subject, body)); err != nil {
		return fmt.Errorf("write SMTP message: %w", err)
	}
	if err := writer.Close(); err != nil {
		return fmt.Errorf("finish SMTP message: %w", err)
	}
	if err := client.Quit(); err != nil {
		return fmt.Errorf("quit SMTP: %w", err)
	}
	return nil
}

// encodeMessage 将主题和正文编码为适合 SMTP 传输的 MIME 邮件。
func encodeMessage(from, to, subject, body string) []byte {
	var message bytes.Buffer
	fmt.Fprintf(&message, "From: %s\r\n", (&mail.Address{Address: from}).String())
	fmt.Fprintf(&message, "To: %s\r\n", (&mail.Address{Address: to}).String())
	fmt.Fprintf(&message, "Subject: %s\r\n", mime.QEncoding.Encode("UTF-8", subject))
	fmt.Fprintf(&message, "Date: %s\r\n", time.Now().Format(time.RFC1123Z))
	message.WriteString("MIME-Version: 1.0\r\n")
	message.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	message.WriteString("Content-Transfer-Encoding: quoted-printable\r\n\r\n")
	writer := quotedprintable.NewWriter(&message)
	writer.Write([]byte(body))
	writer.Close()
	return message.Bytes()
}
