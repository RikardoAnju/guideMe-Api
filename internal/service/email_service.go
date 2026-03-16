package service

import (
    "context"
    "fmt"
    "os"

    "github.com/mailersend/mailersend-go"
)

type EmailService struct {
    client *mailersend.Mailersend
    from   mailersend.From
    appURL string
}

func NewEmailService() *EmailService {
    ms := mailersend.NewMailersend(os.Getenv("MAILERSEND_API_KEY"))

    return &EmailService{
        client: ms,
        from: mailersend.From{
            Name:  os.Getenv("MAILERSEND_FROM_NAME"),
            Email: os.Getenv("MAILERSEND_FROM_EMAIL"),
        },
        appURL: os.Getenv("APP_URL"),
    }
}

func (e *EmailService) SendVerificationEmail(toEmail, toName, token string) error {
    ctx := context.Background()

    verificationLink := fmt.Sprintf("%s/api/auth/verify-email?token=%s", e.appURL, token)

    recipients := []mailersend.Recipient{
        {
            Name:  toName,
            Email: toEmail,
        },
    }
    personalization := []mailersend.Personalization{
        {
            Email: toEmail,
            Data: map[string]interface{}{
                "name":              toName,
                "verification_link": verificationLink,
            },
        },
    }

    message := e.client.Email.NewMessage()
    message.SetFrom(e.from)
    message.SetRecipients(recipients)
    message.SetSubject("Verify Your Email Address")
    message.SetPersonalization(personalization)

    htmlContent := fmt.Sprintf(`
        <h2>Hello, %s!</h2>
        <p>Please verify your email by clicking the button below:</p>
        <a href="%s" style="
            background-color:#4F46E5;
            color:white;
            padding:12px 24px;
            text-decoration:none;
            border-radius:6px;
            display:inline-block;
        ">Verify Email</a>
        <p>Or copy this link: %s</p>
        <p>This link expires in 24 hours.</p>
    `, toName, verificationLink, verificationLink)

    message.SetHTML(htmlContent)

    _, err := e.client.Email.Send(ctx, message)
    return err
}