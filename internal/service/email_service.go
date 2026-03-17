package service

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/mailersend/mailersend-go"
)

type EmailService interface {
	SendVerificationEmail(toEmail, toName, token string) error
}

type mailerSendService struct {
	client    *mailersend.Mailersend
	fromEmail string
	fromName  string
	appURL    string
}

func NewEmailService() EmailService {
	return &mailerSendService{
		client:    mailersend.NewMailersend(os.Getenv("MAILERSEND_API_KEY")),
		fromEmail: os.Getenv("MAILERSEND_FROM_EMAIL"),
		fromName:  os.Getenv("MAILERSEND_FROM_NAME"),
		appURL:    os.Getenv("APP_URL"),
	}
}

func (s *mailerSendService) SendVerificationEmail(toEmail, toName, token string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	verificationLink := fmt.Sprintf("%s/api/v1/auth/verify-email?token=%s", s.appURL, token)

	// Debug log
	fmt.Println("=== SENDING EMAIL ===")
	fmt.Println("To:", toEmail)
	fmt.Println("From:", s.fromEmail)
	fmt.Println("APP_URL:", s.appURL)
	fmt.Println("Verification Link:", verificationLink)
	fmt.Println("API Key:", s.client)

	from := mailersend.From{
		Name:  s.fromName,
		Email: s.fromEmail,
	}
	recipients := []mailersend.Recipient{
		{Name: toName, Email: toEmail},
	}
	htmlContent := fmt.Sprintf(`
        <div style="font-family:sans-serif;max-width:600px;margin:auto;border:1px solid #eee;padding:20px;">
            <h2 style="color: #333;">Verify Your Email</h2>
            <p>Hi %s, please click the button below to verify your account:</p>
            <a href="%s" style="background:#4F46E5;color:white;padding:10px 20px;text-decoration:none;border-radius:5px;display:inline-block;">Verify Now</a>
            <p style="margin-top:20px;font-size:12px;color:#888;">This link expires in 24 hours.</p>
        </div>`, toName, verificationLink)

	message := s.client.Email.NewMessage()
	message.SetFrom(from)
	message.SetRecipients(recipients)
	message.SetSubject("Verify Your Email Address")
	message.SetHTML(htmlContent)

	resp, err := s.client.Email.Send(ctx, message)
	if err != nil {
		fmt.Println("MAILERSEND ERROR:", err)
		fmt.Printf("MAILERSEND RESPONSE: %+v\n", resp)
		return err
	}

	fmt.Println("EMAIL SENT SUCCESS")
	fmt.Printf("RESPONSE: %+v\n", resp)

	return nil
}