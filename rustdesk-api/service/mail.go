package service

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/smtp"
	"strings"
	"time"

	"github.com/lejianwen/rustdesk-api/v2/config"
)

type MailService struct{}

func (ms *MailService) IsConfigured() bool {
	return Config != nil &&
		strings.TrimSpace(Config.Mail.Host) != "" &&
		Config.Mail.Port > 0 &&
		strings.TrimSpace(Config.Mail.From) != ""
}

func normalizeMailConfig(cfg config.Mail) config.Mail {
	cfg.Host = strings.TrimSpace(cfg.Host)
	cfg.Host = strings.TrimPrefix(cfg.Host, "ssl://")
	cfg.Host = strings.TrimPrefix(cfg.Host, "smtp://")
	cfg.Username = strings.TrimSpace(cfg.Username)
	cfg.From = strings.TrimSpace(cfg.From)
	cfg.FromName = strings.TrimSpace(cfg.FromName)
	if cfg.Username == "" {
		cfg.Username = cfg.From
	}
	return cfg
}

func (ms *MailService) Send(to, subject, body string) error {
	if Config == nil {
		return errors.New("mail service not configured")
	}
	return ms.SendWithConfig(Config.Mail, to, subject, body)
}

func (ms *MailService) SendWithConfig(rawCfg config.Mail, to, subject, body string) error {
	cfg := normalizeMailConfig(rawCfg)
	if cfg.Host == "" || cfg.Port <= 0 || cfg.From == "" {
		return errors.New("mail service not configured")
	}
	if strings.TrimSpace(to) == "" {
		return errors.New("recipient email is empty")
	}

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	tlsConfig := &tls.Config{
		ServerName:         cfg.Host,
		InsecureSkipVerify: cfg.SkipVerify,
	}

	var client *smtp.Client
	var err error

	if cfg.UseSSL {
		conn, dialErr := tls.Dial("tcp", addr, tlsConfig)
		if dialErr != nil {
			return dialErr
		}
		client, err = smtp.NewClient(conn, cfg.Host)
	} else {
		client, err = smtp.Dial(addr)
		if err == nil {
			if ok, _ := client.Extension("STARTTLS"); ok {
				if startErr := client.StartTLS(tlsConfig); startErr != nil {
					client.Close()
					return startErr
				}
			}
		}
	}
	if err != nil {
		return err
	}
	defer client.Close()

	if cfg.Username != "" {
		auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)
		if ok, _ := client.Extension("AUTH"); ok {
			if err := client.Auth(auth); err != nil {
				return err
			}
		}
	}

	fromAddress := cfg.From
	if err := client.Mail(fromAddress); err != nil {
		return err
	}
	if err := client.Rcpt(strings.TrimSpace(to)); err != nil {
		return err
	}

	w, err := client.Data()
	if err != nil {
		return err
	}

	from := mail.Address{Name: cfg.FromName, Address: fromAddress}
	headers := []string{
		fmt.Sprintf("From: %s", from.String()),
		fmt.Sprintf("To: %s", strings.TrimSpace(to)),
		fmt.Sprintf("Subject: %s", subject),
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
		fmt.Sprintf("Date: %s", time.Now().Format(time.RFC1123Z)),
	}
	message := strings.Join(headers, "\r\n") + "\r\n\r\n" + body

	if _, err := w.Write([]byte(message)); err != nil {
		_ = w.Close()
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}

	return client.Quit()
}

func (ms *MailService) SendRegisterSuccess(to, username string) error {
	subject := "Registration successful"
	body := fmt.Sprintf("Your account has been created successfully.\n\nUsername: %s\n", username)
	return ms.Send(to, subject, body)
}

func (ms *MailService) ServerReachable() error {
	if !ms.IsConfigured() {
		return errors.New("mail service not configured")
	}
	cfg := normalizeMailConfig(Config.Mail)
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		return err
	}
	return conn.Close()
}
