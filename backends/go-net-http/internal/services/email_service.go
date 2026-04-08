package services

import (
	"context"
	"fmt"
	"log"
	"net/smtp"
	"strings"

	"stacks-base/backends/go-net-http/internal/config"
	"stacks-base/backends/go-net-http/internal/repositories"
)

type EmailSender interface {
	SendRegistrationNotice(ctx context.Context, user repositories.User) error
	SendPasswordResetNotice(ctx context.Context, user repositories.User, resetURL string) error
}

type SMTPEmailService struct {
	host        string
	port        string
	username    string
	password    string
	fromName    string
	fromAddress string
	frontendURL string
}

func NewSMTPEmailService(cfg config.Config) *SMTPEmailService {
	return &SMTPEmailService{
		host:        cfg.SMTPHost,
		port:        cfg.SMTPPort,
		username:    cfg.SMTPUsername,
		password:    cfg.SMTPPassword,
		fromName:    cfg.SMTPFromName,
		fromAddress: cfg.SMTPFromAddress,
		frontendURL: cfg.FrontendBaseURL,
	}
}

func (s *SMTPEmailService) SendRegistrationNotice(_ context.Context, user repositories.User) error {
	auth := smtp.PlainAuth("", s.username, s.password, s.host)
	recipient := []string{user.Email}
	body := strings.Join([]string{
		fmt.Sprintf("From: %s <%s>", s.fromName, s.fromAddress),
		fmt.Sprintf("To: %s", user.Email),
		"Subject: Sua conta no Stacks Base foi criada",
		"MIME-version: 1.0;",
		"Content-Type: text/plain; charset=\"UTF-8\";",
		"",
		fmt.Sprintf("Ola, %s!", user.Name),
		"",
		"Sua conta de referencia no Stacks Base foi criada com sucesso.",
		fmt.Sprintf("Voce pode continuar o fluxo em %s/app.", s.frontendURL),
		"",
		"Esta mensagem foi enviada pelo ambiente de baseline via Mailtrap SMTP.",
	}, "\r\n")

	address := fmt.Sprintf("%s:%s", s.host, s.port)
	if err := smtp.SendMail(address, auth, s.fromAddress, recipient, []byte(body)); err != nil {
		return fmt.Errorf("send registration email: %w", err)
	}

	return nil
}

func (s *SMTPEmailService) SendPasswordResetNotice(_ context.Context, user repositories.User, resetURL string) error {
	auth := smtp.PlainAuth("", s.username, s.password, s.host)
	recipient := []string{user.Email}
	body := strings.Join([]string{
		fmt.Sprintf("From: %s <%s>", s.fromName, s.fromAddress),
		fmt.Sprintf("To: %s", user.Email),
		"Subject: Recuperacao de acesso no Stacks Base",
		"MIME-version: 1.0;",
		"Content-Type: text/plain; charset=\"UTF-8\";",
		"",
		fmt.Sprintf("Ola, %s!", user.Name),
		"",
		"Recebemos uma solicitacao para redefinir a sua senha.",
		fmt.Sprintf("Abra o link abaixo para continuar: %s", resetURL),
		"",
		"Se voce nao solicitou esta acao, ignore este e-mail.",
	}, "\r\n")

	address := fmt.Sprintf("%s:%s", s.host, s.port)
	if err := smtp.SendMail(address, auth, s.fromAddress, recipient, []byte(body)); err != nil {
		return fmt.Errorf("send password reset email: %w", err)
	}

	return nil
}

type LogOnlyEmailService struct{}

func (s LogOnlyEmailService) SendRegistrationNotice(_ context.Context, user repositories.User) error {
	log.Printf("email delivery skipped for %s", user.Email)
	return nil
}

func (s LogOnlyEmailService) SendPasswordResetNotice(_ context.Context, user repositories.User, resetURL string) error {
	log.Printf("password reset email skipped for %s (%s)", user.Email, resetURL)
	return nil
}
