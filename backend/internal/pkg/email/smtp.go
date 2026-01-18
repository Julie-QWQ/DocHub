package email

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"
)

// SMTPConfig SMTP配置
type SMTPConfig struct {
	Host     string // SMTP服务器地址
	Port     int    // SMTP端口
	Username string // 发件人邮箱
	Password string // 邮箱密码或授权码
	From     string // 发件人名称
}

// SMTPClient SMTP客户端
type SMTPClient struct {
	config *SMTPConfig
	auth   smtp.Auth
}

// NewSMTPClient 创建SMTP客户端
func NewSMTPClient(config *SMTPConfig) *SMTPClient {
	auth := smtp.PlainAuth("", config.Username, config.Password, config.Host)
	return &SMTPClient{
		config: config,
		auth:   auth,
	}
}

// SendVerificationCode 发送验证码邮件
func (c *SMTPClient) SendVerificationCode(toEmail, code string, purpose string) error {
	// 根据用途确定邮件主题和内容
	var subject, body string
	switch purpose {
	case "register":
		subject = "UPC-DocHub 注册验证码"
		body = fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
</head>
<body style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
    <div style="background-color: #f8f9fa; padding: 30px; border-radius: 5px;">
        <h2 style="color: #333;">欢迎注册 UPC-DocHub</h2>
        <p>您好,</p>
        <p>感谢您注册 UPC-DocHub 学习资料托管平台。您的验证码是:</p>
        <div style="background-color: #007bff; color: white; padding: 15px; text-align: center; font-size: 24px; font-weight: bold; border-radius: 5px; margin: 20px 0;">
            %s
        </div>
        <p>验证码有效期为 <strong>10分钟</strong>,请尽快完成验证。</p>
        <p style="color: #666; font-size: 12px;">如果这不是您的操作,请忽略此邮件。</p>
    </div>
</body>
</html>`, code)
	case "login":
		subject = "UPC-DocHub 登录验证码"
		body = fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
</head>
<body style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
    <div style="background-color: #f8f9fa; padding: 30px; border-radius: 5px;">
        <h2 style="color: #333;">登录验证码</h2>
        <p>您好,</p>
        <p>您正在登录 Julie 的 UPC-DocHub 学习资料托管平台。您的验证码是:</p>
        <div style="background-color: #28a745; color: white; padding: 15px; text-align: center; font-size: 24px; font-weight: bold; border-radius: 5px; margin: 20px 0;">
            %s
        </div>
        <p>验证码有效期为 <strong>10分钟</strong>,请尽快完成验证。</p>
        <p style="color: #666; font-size: 12px;">如果这不是您的操作,请立即修改密码。</p>
    </div>
</body>
</html>`, code)
	case "reset_password":
		subject = "UPC-DocHub 重置密码验证码"
		body = fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
</head>
<body style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
    <div style="background-color: #f8f9fa; padding: 30px; border-radius: 5px;">
        <h2 style="color: #333;">重置密码验证码</h2>
        <p>您好,</p>
        <p>您正在重置 UPC-DocHub 账号密码。您的验证码是:</p>
        <div style="background-color: #dc3545; color: white; padding: 15px; text-align: center; font-size: 24px; font-weight: bold; border-radius: 5px; margin: 20px 0;">
            %s
        </div>
        <p>验证码有效期为 <strong>10分钟</strong>,请尽快完成验证。</p>
        <p style="color: #666; font-size: 12px;">如果这不是您的操作,请忽略此邮件。</p>
    </div>
</body>
</html>`, code)
	default:
		return fmt.Errorf("未知的邮件用途: %s", purpose)
	}

	// 构建邮件内容
	message := fmt.Sprintf("From: %s\r\n", c.config.Username)
	message += fmt.Sprintf("To: %s\r\n", toEmail)
	message += fmt.Sprintf("Subject: %s\r\n", subject)
	message += "MIME-Version: 1.0\r\n"
	message += "Content-Type: text/html; charset=UTF-8\r\n\r\n"
	message += body

	// 发送邮件
	addr := fmt.Sprintf("%s:%d", c.config.Host, c.config.Port)

	// 使用TLS加密连接
	tlsconfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         c.config.Host,
	}

	conn, err := tls.Dial("tcp", addr, tlsconfig)
	if err != nil {
		return fmt.Errorf("连接SMTP服务器失败: %w", err)
	}

	client, err := smtp.NewClient(conn, c.config.Host)
	if err != nil {
		return fmt.Errorf("创建SMTP客户端失败: %w", err)
	}
	defer client.Close()

	// 认证
	if err := client.Auth(c.auth); err != nil {
		return fmt.Errorf("SMTP认证失败: %w", err)
	}

	// 设置发件人
	if err := client.Mail(c.config.Username); err != nil {
		return fmt.Errorf("设置发件人失败: %w", err)
	}

	// 设置收件人
	if err := client.Rcpt(toEmail); err != nil {
		return fmt.Errorf("设置收件人失败: %w", err)
	}

	// 发送邮件内容
	wc, err := client.Data()
	if err != nil {
		return fmt.Errorf("获取数据写入器失败: %w", err)
	}
	defer wc.Close()

	_, err = wc.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("写入邮件内容失败: %w", err)
	}

	return nil
}

// SendGenericEmail 发送普通邮件
func (c *SMTPClient) SendGenericEmail(toEmail, subject, htmlBody string) error {
	message := fmt.Sprintf("From: %s\r\n", c.config.Username)
	message += fmt.Sprintf("To: %s\r\n", toEmail)
	message += fmt.Sprintf("Subject: %s\r\n", subject)
	message += "MIME-Version: 1.0\r\n"
	message += "Content-Type: text/html; charset=UTF-8\r\n\r\n"
	message += htmlBody

	addr := fmt.Sprintf("%s:%d", c.config.Host, c.config.Port)

	tlsconfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         c.config.Host,
	}

	conn, err := tls.Dial("tcp", addr, tlsconfig)
	if err != nil {
		return fmt.Errorf("连接SMTP服务器失败: %w", err)
	}

	client, err := smtp.NewClient(conn, c.config.Host)
	if err != nil {
		return fmt.Errorf("创建SMTP客户端失败: %w", err)
	}
	defer client.Close()

	if err := client.Auth(c.auth); err != nil {
		return fmt.Errorf("SMTP认证失败: %w", err)
	}

	if err := client.Mail(c.config.Username); err != nil {
		return fmt.Errorf("设置发件人失败: %w", err)
	}

	if err := client.Rcpt(toEmail); err != nil {
		return fmt.Errorf("设置收件人失败: %w", err)
	}

	wc, err := client.Data()
	if err != nil {
		return fmt.Errorf("获取数据写入器失败: %w", err)
	}
	defer wc.Close()

	_, err = wc.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("写入邮件内容失败: %w", err)
	}

	return nil
}

// ValidateEmailFormat 验证邮箱格式
func ValidateEmailFormat(email string) bool {
	email = strings.TrimSpace(email)
	if len(email) < 3 || len(email) > 254 {
		return false
	}

	atIndex := strings.LastIndex(email, "@")
	if atIndex <= 0 || atIndex == len(email)-1 {
		return false
	}

	localPart := email[:atIndex]
	domainPart := email[atIndex+1:]

	if len(localPart) == 0 || len(domainPart) < 3 {
		return false
	}

	if !strings.Contains(domainPart, ".") {
		return false
	}

	return true
}
