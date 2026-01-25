package email

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"

	"github.com/erp-cosmetics/notification-service/internal/domain/entity"
	"go.uber.org/zap"
)

// SMTPConfig holds SMTP configuration
type SMTPConfig struct {
	Host      string
	Port      int
	Username  string
	Password  string
	FromEmail string
	FromName  string
	UseTLS    bool
}

// Sender defines email sender interface
type Sender interface {
	Send(to, subject, bodyHTML, bodyText string) (*entity.EmailLog, error)
	SendWithCC(to []string, cc []string, subject, bodyHTML, bodyText string) (*entity.EmailLog, error)
}

type smtpSender struct {
	config *SMTPConfig
	logger *zap.Logger
}

// NewSMTPSender creates a new SMTP email sender
func NewSMTPSender(config *SMTPConfig, logger *zap.Logger) Sender {
	return &smtpSender{
		config: config,
		logger: logger,
	}
}

func (s *smtpSender) Send(to, subject, bodyHTML, bodyText string) (*entity.EmailLog, error) {
	return s.SendWithCC([]string{to}, nil, subject, bodyHTML, bodyText)
}

func (s *smtpSender) SendWithCC(to []string, cc []string, subject, bodyHTML, bodyText string) (*entity.EmailLog, error) {
	// Create email log
	emailLog := &entity.EmailLog{
		FromEmail: s.config.FromEmail,
		FromName:  s.config.FromName,
		ToEmail:   strings.Join(to, ","),
		Subject:   subject,
		BodyHTML:  bodyHTML,
		BodyText:  bodyText,
		Status:    entity.EmailStatusPending,
	}

	// Build email message
	msg := s.buildMessage(to, cc, subject, bodyHTML, bodyText)

	// Get all recipients
	recipients := append(to, cc...)

	// Send email
	var err error
	if s.config.UseTLS {
		err = s.sendWithTLS(recipients, msg)
	} else {
		err = s.sendWithStartTLS(recipients, msg)
	}

	if err != nil {
		s.logger.Error("Failed to send email",
			zap.Strings("to", to),
			zap.String("subject", subject),
			zap.Error(err),
		)
		emailLog.MarkAsFailed("SMTP_ERROR", err.Error())
		return emailLog, err
	}

	emailLog.MarkAsSent("", "OK")
	s.logger.Info("Email sent successfully",
		zap.Strings("to", to),
		zap.String("subject", subject),
	)

	return emailLog, nil
}

func (s *smtpSender) buildMessage(to, cc []string, subject, bodyHTML, bodyText string) []byte {
	var msg strings.Builder

	// Headers
	msg.WriteString(fmt.Sprintf("From: %s <%s>\r\n", s.config.FromName, s.config.FromEmail))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(to, ", ")))
	if len(cc) > 0 {
		msg.WriteString(fmt.Sprintf("Cc: %s\r\n", strings.Join(cc, ", ")))
	}
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	msg.WriteString("MIME-Version: 1.0\r\n")

	// Check if we have both HTML and text
	if bodyHTML != "" && bodyText != "" {
		boundary := "boundary-erp-notification"
		msg.WriteString(fmt.Sprintf("Content-Type: multipart/alternative; boundary=%s\r\n\r\n", boundary))

		// Plain text part
		msg.WriteString(fmt.Sprintf("--%s\r\n", boundary))
		msg.WriteString("Content-Type: text/plain; charset=UTF-8\r\n\r\n")
		msg.WriteString(bodyText)
		msg.WriteString("\r\n")

		// HTML part
		msg.WriteString(fmt.Sprintf("--%s\r\n", boundary))
		msg.WriteString("Content-Type: text/html; charset=UTF-8\r\n\r\n")
		msg.WriteString(bodyHTML)
		msg.WriteString("\r\n")

		msg.WriteString(fmt.Sprintf("--%s--\r\n", boundary))
	} else if bodyHTML != "" {
		msg.WriteString("Content-Type: text/html; charset=UTF-8\r\n\r\n")
		msg.WriteString(bodyHTML)
	} else {
		msg.WriteString("Content-Type: text/plain; charset=UTF-8\r\n\r\n")
		msg.WriteString(bodyText)
	}

	return []byte(msg.String())
}

func (s *smtpSender) sendWithTLS(recipients []string, msg []byte) error {
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

	tlsConfig := &tls.Config{
		ServerName: s.config.Host,
	}

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, s.config.Host)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	defer client.Close()

	return s.sendViaClient(client, recipients, msg)
}

func (s *smtpSender) sendWithStartTLS(recipients []string, msg []byte) error {
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

	client, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("failed to dial: %w", err)
	}
	defer client.Close()

	// Try STARTTLS if available
	if ok, _ := client.Extension("STARTTLS"); ok {
		tlsConfig := &tls.Config{
			ServerName: s.config.Host,
		}
		if err := client.StartTLS(tlsConfig); err != nil {
			return fmt.Errorf("failed to start TLS: %w", err)
		}
	}

	return s.sendViaClient(client, recipients, msg)
}

func (s *smtpSender) sendViaClient(client *smtp.Client, recipients []string, msg []byte) error {
	// Authenticate if credentials provided
	if s.config.Username != "" && s.config.Password != "" {
		auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)
		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}
	}

	// Set sender
	if err := client.Mail(s.config.FromEmail); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	// Set recipients
	for _, recipient := range recipients {
		if recipient == "" {
			continue
		}
		if err := client.Rcpt(recipient); err != nil {
			return fmt.Errorf("failed to add recipient %s: %w", recipient, err)
		}
	}

	// Send message body
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %w", err)
	}

	if _, err := w.Write(msg); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	if err := w.Close(); err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	return client.Quit()
}

// MockSender is a mock email sender for testing
type MockSender struct {
	SentEmails []*entity.EmailLog
	logger     *zap.Logger
}

// NewMockSender creates a new mock email sender
func NewMockSender(logger *zap.Logger) *MockSender {
	return &MockSender{
		SentEmails: make([]*entity.EmailLog, 0),
		logger:     logger,
	}
}

func (m *MockSender) Send(to, subject, bodyHTML, bodyText string) (*entity.EmailLog, error) {
	return m.SendWithCC([]string{to}, nil, subject, bodyHTML, bodyText)
}

func (m *MockSender) SendWithCC(to []string, cc []string, subject, bodyHTML, bodyText string) (*entity.EmailLog, error) {
	emailLog := &entity.EmailLog{
		FromEmail: "noreply@example.com",
		FromName:  "ERP System",
		ToEmail:   strings.Join(to, ","),
		Subject:   subject,
		BodyHTML:  bodyHTML,
		BodyText:  bodyText,
	}
	emailLog.MarkAsSent("mock-message-id", "OK")
	
	m.SentEmails = append(m.SentEmails, emailLog)
	m.logger.Info("Mock email sent",
		zap.Strings("to", to),
		zap.String("subject", subject),
	)
	
	return emailLog, nil
}
