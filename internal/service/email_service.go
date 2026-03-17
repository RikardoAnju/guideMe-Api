package service

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/mailersend/mailersend-go"
)

// Tambah SendOTPEmail ke interface
type EmailService interface {
	SendVerificationEmail(toEmail, toName, token string) error
	SendOTPEmail(toEmail, toName, otp string) error
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

	fmt.Println("=== SENDING VERIFICATION EMAIL ===")
	fmt.Println("To:", toEmail)
	fmt.Println("Verification Link:", verificationLink)

	from := mailersend.From{
		Name:  s.fromName,
		Email: s.fromEmail,
	}

	recipients := []mailersend.Recipient{
		{Name: toName, Email: toEmail},
	}

	htmlContent := fmt.Sprintf(`
		<div style="font-family:sans-serif;max-width:600px;margin:auto;border:1px solid #eee;padding:20px;">
			<h2 style="color:#333;">Verify Your Email</h2>
			<p>Hi %s, please click the button below to verify your account:</p>
			<a href="%s" style="background:#4F46E5;color:white;padding:10px 20px;text-decoration:none;border-radius:5px;display:inline-block;">
				Verify Now
			</a>
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

	fmt.Println("VERIFICATION EMAIL SENT SUCCESS")
	return nil
}

// SendOTPEmail — kirim kode OTP untuk reset password
func (s *mailerSendService) SendOTPEmail(toEmail, toName, otp string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	fmt.Println("=== SENDING OTP EMAIL ===")
	fmt.Println("To:", toEmail)
	fmt.Println("OTP:", otp)

	from := mailersend.From{
		Name:  s.fromName,
		Email: s.fromEmail,
	}

	recipients := []mailersend.Recipient{
		{Name: toName, Email: toEmail},
	}

	htmlContent := fmt.Sprintf(`
		<div style="font-family:sans-serif;max-width:600px;margin:auto;border:1px solid #eee;padding:20px;">
			<h2 style="color:#333;">Reset Password</h2>
			<p>Hi %s, gunakan kode OTP berikut untuk mereset password kamu:</p>
			<div style="font-size:36px;font-weight:bold;letter-spacing:12px;text-align:center;
						padding:24px;background:#f4f4f4;border-radius:8px;margin:24px 0;
						color:#4F46E5;">
				%s
			</div>
			<p>Kode ini berlaku selama <strong>10 menit</strong>.</p>
			<p style="font-size:12px;color:#888;">
				Jika kamu tidak merasa melakukan permintaan ini, abaikan email ini.
			</p>
		</div>`, toName, otp)

	message := s.client.Email.NewMessage()
	message.SetFrom(from)
	message.SetRecipients(recipients)
	message.SetSubject("Kode OTP Reset Password")
	message.SetHTML(htmlContent)

	resp, err := s.client.Email.Send(ctx, message)
	if err != nil {
		fmt.Println("MAILERSEND OTP ERROR:", err)
		fmt.Printf("MAILERSEND RESPONSE: %+v\n", resp)
		return err
	}

	fmt.Println("OTP EMAIL SENT SUCCESS")
	return nil
}