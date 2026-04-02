package services

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	"strings"

	"bbsgo/config"
)

// EmailService 邮件服务结构
// 提供邮件发送功能，支持 SMTP SSL/TLS 两种方式
type EmailService struct {
	enabled  bool   // 是否启用邮件服务
	host     string // SMTP 服务器地址
	port     int    // SMTP 端口
	user     string // 用户名
	password string // 密码
	from     string // 发件人地址
	fromName string // 发件人名称
}

// NewEmailService 创建邮件服务实例
// 从配置中读取邮件相关参数
func NewEmailService() *EmailService {
	return &EmailService{
		enabled:  config.GetConfigBool("email_enabled", false),
		host:     config.GetConfig("email_host"),
		port:     config.GetConfigInt("email_port", 465),
		user:     config.GetConfig("email_user"),
		password: config.GetConfig("email_password"),
		from:     config.GetConfig("email_from"),
		fromName: config.GetConfig("email_from_name"),
	}
}

// Send 发送邮件
// to: 收件人地址
// subject: 邮件主题
// body: 邮件正文（HTML）
// 返回: 错误信息
func (s *EmailService) Send(to, subject, body string) error {
	// 检查邮件服务是否启用
	if !s.enabled {
		return fmt.Errorf("email service is disabled")
	}

	log.Printf("[email] sending email, to: %s, subject: %s", to, subject)
	log.Printf("[email] config - host: %s, port: %d, user: %s, from: %s, fromName: %s",
		s.host, s.port, s.user, s.from, s.fromName)

	addr := fmt.Sprintf("%s:%d", s.host, s.port)

	// 根据端口选择发送方式
	if s.port == 465 {
		return s.sendMailWithSSL(addr, s.user, s.password, s.from, to, subject, body)
	}
	return s.sendMailWithTLS(addr, s.user, s.password, s.from, to, subject, body)
}

// sendMailWithTLS 使用 TLS 方式发送邮件
func (s *EmailService) sendMailWithTLS(addr, user, password, from, to, subject, body string) error {
	// 创建 SMTP 认证
	auth := smtp.PlainAuth("", user, password, strings.Split(addr, ":")[0])

	// 建立 TLS 连接
	conn, err := tls.Dial("tcp", addr, &tls.Config{ServerName: strings.Split(addr, ":")[0]})
	if err != nil {
		log.Printf("[email] failed to connect to SMTP server (TLS): %v", err)
		return fmt.Errorf("failed to connect to SMTP server: %v", err)
	}
	defer conn.Close()

	// 创建 SMTP 客户端
	c, err := smtp.NewClient(conn, strings.Split(addr, ":")[0])
	if err != nil {
		log.Printf("[email] failed to create SMTP client (TLS): %v", err)
		return fmt.Errorf("failed to create SMTP client: %v", err)
	}
	defer c.Quit()

	// 执行 AUTH 认证
	if ok, _ := c.Extension("AUTH"); ok {
		if err := c.Auth(auth); err != nil {
			log.Printf("[email] failed to authenticate (TLS): %v", err)
			return fmt.Errorf("failed to authenticate: %v", err)
		}
	}

	// 设置发件人
	if err := c.Mail(from); err != nil {
		log.Printf("[email] failed to set sender (TLS): %v", err)
		return fmt.Errorf("failed to set sender: %v", err)
	}

	// 设置收件人
	if err := c.Rcpt(to); err != nil {
		log.Printf("[email] failed to set recipient (TLS): %v", err)
		return fmt.Errorf("failed to set recipient: %v", err)
	}

	// 获取数据写入器
	w, err := c.Data()
	if err != nil {
		log.Printf("[email] failed to get data writer (TLS): %v", err)
		return fmt.Errorf("failed to get data writer: %v", err)
	}

	// 构建并写入邮件内容
	msg := s.buildMessage(from, to, subject, body)
	if _, err := w.Write([]byte(msg)); err != nil {
		log.Printf("[email] failed to write message (TLS): %v", err)
		return fmt.Errorf("failed to write message: %v", err)
	}

	if err := w.Close(); err != nil {
		log.Printf("[email] failed to close data writer (TLS): %v", err)
		return fmt.Errorf("failed to close data writer: %v", err)
	}

	log.Printf("[email] email sent successfully (TLS), to: %s", to)
	return nil
}

// sendMailWithSSL 使用 SSL 方式发送邮件（适用于465端口）
func (s *EmailService) sendMailWithSSL(addr, user, password, from, to, subject, body string) error {
	// 创建 SMTP 认证
	auth := smtp.PlainAuth("", user, password, strings.Split(addr, ":")[0])

	// 建立 SSL 连接
	conn, err := tls.Dial("tcp", addr, &tls.Config{
		ServerName:         strings.Split(addr, ":")[0],
		InsecureSkipVerify: true, // 跳过证书验证
	})
	if err != nil {
		log.Printf("[email] failed to connect to SMTP server (SSL): %v", err)
		return fmt.Errorf("failed to connect to SMTP server: %v", err)
	}
	defer conn.Close()

	// 创建 SMTP 客户端
	c, err := smtp.NewClient(conn, strings.Split(addr, ":")[0])
	if err != nil {
		log.Printf("[email] failed to create SMTP client (SSL): %v", err)
		return fmt.Errorf("failed to create SMTP client: %v", err)
	}
	defer c.Quit()

	// 执行 AUTH 认证
	if ok, _ := c.Extension("AUTH"); ok {
		if err := c.Auth(auth); err != nil {
			log.Printf("[email] failed to authenticate (SSL): %v", err)
			return fmt.Errorf("failed to authenticate: %v", err)
		}
	}

	// 设置发件人
	if err := c.Mail(from); err != nil {
		log.Printf("[email] failed to set sender (SSL): %v", err)
		return fmt.Errorf("failed to set sender: %v", err)
	}

	// 设置收件人
	if err := c.Rcpt(to); err != nil {
		log.Printf("[email] failed to set recipient (SSL): %v", err)
		return fmt.Errorf("failed to set recipient: %v", err)
	}

	// 获取数据写入器
	w, err := c.Data()
	if err != nil {
		log.Printf("[email] failed to get data writer (SSL): %v", err)
		return fmt.Errorf("failed to get data writer: %v", err)
	}

	// 构建并写入邮件内容
	msg := s.buildMessage(from, to, subject, body)
	if _, err := w.Write([]byte(msg)); err != nil {
		log.Printf("[email] failed to write message (SSL): %v", err)
		return fmt.Errorf("failed to write message: %v", err)
	}

	if err := w.Close(); err != nil {
		log.Printf("[email] failed to close data writer (SSL): %v", err)
		return fmt.Errorf("failed to close data writer: %v", err)
	}

	log.Printf("[email] email sent successfully (SSL), to: %s", to)
	return nil
}

// buildMessage 构建邮件内容
// 返回标准 MIME 格式的邮件内容
func (s *EmailService) buildMessage(from, to, subject, body string) string {
	msg := fmt.Sprintf("From: %s <%s>\r\n", s.fromName, from)
	msg += fmt.Sprintf("To: %s\r\n", to)
	msg += fmt.Sprintf("Subject: %s\r\n", subject)
	msg += "MIME-Version: 1.0\r\n"
	msg += "Content-Type: text/html; charset=\"UTF-8\"\r\n\r\n"
	msg += body
	return msg
}

// SendVerificationCode 发送验证码邮件
// to: 收件人邮箱
// code: 验证码
// 返回: 错误信息
func SendVerificationCode(to, code string) error {
	service := NewEmailService()
	if !service.enabled {
		return fmt.Errorf("email service is disabled")
	}

	subject := "验证码 - " + service.fromName
	body := fmt.Sprintf(`
		<html>
		<body>
			<h2>您的验证码</h2>
			<p>您的验证码是：<strong style="font-size: 24px; color: #4F46E5;">%s</strong></p>
			<p>验证码有效期为5分钟，请尽快使用。</p>
			<p>如果这不是您的操作，请忽略此邮件。</p>
			<br>
			<p>%s</p>
		</body>
		</html>
	`, code, service.fromName)

	return service.Send(to, subject, body)
}
