# Email

This guide covers email configuration, SMTP setup, and sending emails in Laju Go.

## Overview

Laju Go includes a mailer service for sending emails with the following features:

- **SMTP Support** - Send emails via SMTP
- **HTML Templates** - Rich HTML email templates
- **Password Reset** - Built-in password reset emails
- **Queue Support** - Async email sending (optional)

## Configuration

### Environment Variables

```bash
# .env
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASS=your-app-password
FROM_EMAIL=noreply@example.com
FROM_NAME=Laju Go
```

### SMTP Settings by Provider

#### Gmail

```bash
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASS=your-16-character-app-password
FROM_EMAIL=noreply@example.com
FROM_NAME=Your App Name
```

> 🔐 **Important**: Use an [App Password](https://support.google.com/accounts/answer/185833), not your regular Gmail password.

#### SendGrid

```bash
SMTP_HOST=smtp.sendgrid.com
SMTP_PORT=587
SMTP_USER=apikey
SMTP_PASS=your-sendgrid-api-key
FROM_EMAIL=verified-sender@yourdomain.com
FROM_NAME=Your App Name
```

#### Mailgun

```bash
SMTP_HOST=smtp.mailgun.org
SMTP_PORT=587
SMTP_USER=postmaster@yourdomain.mailgun.org
SMTP_PASS=your-mailgun-api-key
FROM_EMAIL=noreply@yourdomain.com
FROM_NAME=Your App Name
```

#### Office 365

```bash
SMTP_HOST=smtp.office365.com
SMTP_PORT=587
SMTP_USER=your-email@yourdomain.com
SMTP_PASS=your-password
FROM_EMAIL=noreply@yourdomain.com
FROM_NAME=Your App Name
```

## Mailer Service

### Implementation

```go
// app/services/mailer.go
package services

import (
    "bytes"
    "crypto/rand"
    "encoding/base64"
    "fmt"
    "html/template"
    "net/smtp"
    "os"
    "strconv"
    "time"
)

type MailerService struct {
    smtpHost     string
    smtpPort     int
    smtpUser     string
    smtpPass     string
    fromEmail    string
    fromName     string
}

type Email struct {
    To      string
    Subject string
    Body    string
    IsHTML  bool
}

func NewMailerService() *MailerService {
    smtpHost := os.Getenv("SMTP_HOST")
    smtpPort, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
    smtpUser := os.Getenv("SMTP_USER")
    smtpPass := os.Getenv("SMTP_PASS")
    fromEmail := os.Getenv("FROM_EMAIL")
    fromName := os.Getenv("FROM_NAME")
    
    if fromName == "" {
        fromName = "Laju Go"
    }
    
    return &MailerService{
        smtpHost:  smtpHost,
        smtpPort:  smtpPort,
        smtpUser:  smtpUser,
        smtpPass:  smtpPass,
        fromEmail: fromEmail,
        fromName:  fromName,
    }
}

func (m *MailerService) Send(email Email) error {
    // Create message
    var body bytes.Buffer
    
    // Headers
    fmt.Fprintf(&body, "From: %s <%s>\r\n", m.fromName, m.fromEmail)
    fmt.Fprintf(&body, "To: %s\r\n", email.To)
    fmt.Fprintf(&body, "Subject: %s\r\n", email.Subject)
    
    if email.IsHTML {
        fmt.Fprintf(&body, "MIME-version: 1.0\r\n")
        fmt.Fprintf(&body, "Content-Type: text/html; charset=\"UTF-8\"\r\n")
    }
    fmt.Fprintf(&body, "\r\n")
    fmt.Fprintf(&body, "%s", email.Body)
    
    // Send email
    addr := fmt.Sprintf("%s:%d", m.smtpHost, m.smtpPort)
    
    var auth smtp.Auth
    if m.smtpUser != "" && m.smtpPass != "" {
        auth = smtp.PlainAuth("", m.smtpUser, m.smtpPass, m.smtpHost)
    }
    
    return smtp.SendMail(addr, auth, m.fromEmail, []string{email.To}, body.Bytes())
}

func (m *MailerService) SendTemplate(to, subject, templateFile string, data interface{}) error {
    // Parse template
    tmpl, err := template.ParseFiles(templateFile)
    if err != nil {
        return err
    }
    
    // Execute template
    var body bytes.Buffer
    if err := tmpl.Execute(&body, data); err != nil {
        return err
    }
    
    // Send email
    return m.Send(Email{
        To:      to,
        Subject: subject,
        Body:    body.String(),
        IsHTML:  true,
    })
}
```

## Password Reset Email

### Implementation

```go
// app/services/mailer.go
func (m *MailerService) SendPasswordResetEmail(to, token string) error {
    resetURL := fmt.Sprintf("%s/reset-password/%s", os.Getenv("APP_URL"), token)
    
    data := map[string]interface{}{
        "ResetURL": resetURL,
        "AppName":  m.fromName,
    }
    
    return m.SendTemplate(to, "Password Reset Request", "templates/emails/password-reset.html", data)
}
```

### Email Template

```html
<!-- templates/emails/password-reset.html -->
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
        }
        .container {
            background: #f9f9f9;
            border-radius: 8px;
            padding: 30px;
        }
        .button {
            display: inline-block;
            background: #007bff;
            color: white;
            text-decoration: none;
            padding: 12px 24px;
            border-radius: 4px;
            margin-top: 20px;
        }
        .footer {
            margin-top: 30px;
            font-size: 12px;
            color: #666;
            text-align: center;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Password Reset Request</h1>
        
        <p>You requested to reset your password for {{.AppName}}.</p>
        
        <p>Click the button below to reset your password:</p>
        
        <p>
            <a href="{{.ResetURL}}" class="button">Reset Password</a>
        </p>
        
        <p>Or copy and paste this link into your browser:</p>
        <p style="word-break: break-all;">{{.ResetURL}}</p>
        
        <p><strong>This link will expire in 1 hour.</strong></p>
        
        <p>If you didn't request this, you can safely ignore this email.</p>
        
        <div class="footer">
            <p>&copy; {{.AppName}}. All rights reserved.</p>
        </div>
    </div>
</body>
</html>
```

## Usage in Handlers

### Password Reset Handler

```go
// app/handlers/password-reset.go
func (h *PasswordResetHandler) SendResetLink(c *fiber.Ctx) error {
    var req dto.ForgotPasswordRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
    }
    
    // Generate reset token
    token, err := h.userService.GenerateResetToken(req.Email)
    if err != nil {
        // Don't reveal if email exists
        return c.JSON(fiber.Map{
            "message": "If the email exists, a reset link has been sent",
        })
    }
    
    // Send reset email
    err = h.mailerService.SendPasswordResetEmail(req.Email, token)
    if err != nil {
        // Log error but don't expose to user
        log.Printf("Failed to send reset email: %v", err)
    }
    
    return c.JSON(fiber.Map{
        "message": "If the email exists, a reset link has been sent",
    })
}
```

## Welcome Email

### Template

```html
<!-- templates/emails/welcome.html -->
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
        }
        .container {
            background: #f9f9f9;
            border-radius: 8px;
            padding: 30px;
        }
        .button {
            display: inline-block;
            background: #28a745;
            color: white;
            text-decoration: none;
            padding: 12px 24px;
            border-radius: 4px;
            margin-top: 20px;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Welcome to {{.AppName}}!</h1>
        
        <p>Hi {{.Name}},</p>
        
        <p>Thanks for signing up! We're excited to have you on board.</p>
        
        <p>Get started by exploring your dashboard:</p>
        
        <p>
            <a href="{{.DashboardURL}}" class="button">Go to Dashboard</a>
        </p>
        
        <p>If you have any questions, feel free to reach out to our support team.</p>
        
        <p>Best regards,<br>The {{.AppName}} Team</p>
    </div>
</body>
</html>
```

### Send Welcome Email

```go
// app/services/mailer.go
func (m *MailerService) SendWelcomeEmail(to, name string) error {
    data := map[string]interface{}{
        "Name":         name,
        "AppName":      m.fromName,
        "DashboardURL": os.Getenv("APP_URL") + "/app",
    }
    
    return m.SendTemplate(to, "Welcome to "+m.fromName, "templates/emails/welcome.html", data)
}

// Usage in registration handler
func (h *AuthHandler) Register(c *fiber.Ctx) error {
    // ... register user ...
    
    // Send welcome email
    go h.mailerService.SendWelcomeEmail(user.Email, user.Name)
    
    return c.Redirect("/app")
}
```

## Async Email Sending

### Use Goroutines

```go
// Send email asynchronously to avoid blocking
go func() {
    err := mailerService.SendPasswordResetEmail(email, token)
    if err != nil {
        log.Printf("Failed to send email: %v", err)
    }
}()
```

### Email Queue (Advanced)

```go
// app/services/email_queue.go
type EmailQueue struct {
    queue chan Email
}

func NewEmailQueue(size int) *EmailQueue {
    eq := &EmailQueue{
        queue: make(chan Email, size),
    }
    
    // Start worker
    go eq.worker()
    
    return eq
}

func (eq *EmailQueue) worker() {
    for email := range eq.queue {
        // Send email
        if err := mailerService.Send(email); err != nil {
            log.Printf("Failed to send email: %v", err)
        }
    }
}

func (eq *EmailQueue) Send(email Email) {
    eq.queue <- email
}

// Usage
emailQueue := NewEmailQueue(100)
emailQueue.Send(Email{
    To:      "user@example.com",
    Subject: "Welcome",
    Body:    "Hello!",
})
```

## Email Templates

### Base Template

```html
<!-- templates/emails/base.html -->
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        /* Base styles */
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            line-height: 1.6;
            color: #333;
            background: #f5f5f5;
            margin: 0;
            padding: 20px;
        }
        .wrapper {
            max-width: 600px;
            margin: 0 auto;
            background: white;
            border-radius: 8px;
            overflow: hidden;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        .header {
            background: #007bff;
            color: white;
            padding: 20px;
            text-align: center;
        }
        .content {
            padding: 30px;
        }
        .footer {
            background: #f9f9f9;
            padding: 20px;
            text-align: center;
            font-size: 12px;
            color: #666;
        }
        .button {
            display: inline-block;
            background: #007bff;
            color: white;
            text-decoration: none;
            padding: 12px 24px;
            border-radius: 4px;
        }
    </style>
</head>
<body>
    <div class="wrapper">
        <div class="header">
            <h1>{{.AppName}}</h1>
        </div>
        <div class="content">
            {{template "content" .}}
        </div>
        <div class="footer">
            <p>&copy; {{.AppName}}. All rights reserved.</p>
            <p>{{.AppURL}}</p>
        </div>
    </div>
</body>
</html>
```

### Using Base Template

```go
// Parse template with base
tmpl := template.Must(template.ParseFiles(
    "templates/emails/base.html",
    "templates/emails/password-reset.html",
))

// Define content template
// In password-reset.html:
// {{define "content"}}...content...{{end}}
```

## Testing

### Test Email Sending

```go
// services/mailer_test.go
func TestMailerService_Send(t *testing.T) {
    // Skip in CI
    if os.Getenv("CI") != "" {
        t.Skip("Skipping email test in CI")
    }
    
    mailer := NewMailerService()
    
    err := mailer.Send(Email{
        To:      "test@example.com",
        Subject: "Test Email",
        Body:    "This is a test email",
        IsHTML:  false,
    })
    
    if err != nil {
        t.Errorf("Failed to send email: %v", err)
    }
}
```

### Mock Mailer for Testing

```go
// services/mock_mailer.go
type MockMailer struct {
    SentEmails []Email
}

func (m *MockMailer) Send(email Email) error {
    m.SentEmails = append(m.SentEmails, email)
    return nil
}

// Usage in tests
mockMailer := &MockMailer{}
handler := NewAuthHandler(authService, mockMailer, store)

// Trigger email send
// ...

// Verify email was sent
if len(mockMailer.SentEmails) != 1 {
    t.Errorf("Expected 1 email, got %d", len(mockMailer.SentEmails))
}
```

## Best Practices

### 1. Use Environment Variables

```bash
# ✅ Good: Configuration in .env
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587

# ❌ Bad: Hardcoded values
smtpHost := "smtp.gmail.com"
```

### 2. Handle Errors Gracefully

```go
// ✅ Good: Log error but don't block user flow
err = mailer.SendPasswordResetEmail(email, token)
if err != nil {
    log.Printf("Failed to send email: %v", err)
}
return c.JSON(fiber.Map{"message": "Reset link sent"})

// ❌ Bad: Expose error to user
if err != nil {
    return c.Status(500).JSON(fiber.Map{"error": err.Error()})
}
```

### 3. Use Async Sending

```go
// ✅ Good: Non-blocking
go mailer.SendWelcomeEmail(user.Email, user.Name)

// ❌ Bad: Blocks response
mailer.SendWelcomeEmail(user.Email, user.Name)
```

### 4. Validate Email Addresses

```go
import "net/mail"

func isValidEmail(email string) bool {
    _, err := mail.ParseAddress(email)
    return err == nil
}
```

### 5. Use App Passwords

For Gmail and other providers, use app-specific passwords:

```bash
# ✅ Good: App password (16 characters)
SMTP_PASS=abcd-efgh-ijkl-mnop

# ❌ Bad: Regular password (insecure)
SMTP_PASS=my-regular-password
```

## Troubleshooting

### Authentication Failed

**Error**: `smtp: server doesn't support AUTH`

**Solutions**:
1. Use port 587 (TLS) instead of 465 (SSL)
2. Enable "Less secure app access" (not recommended)
3. Use App Password for Gmail

### Connection Timeout

**Error**: `dial tcp: i/o timeout`

**Solutions**:
1. Check firewall settings
2. Verify SMTP host and port
3. Check network connectivity

### Email Not Received

**Solutions**:
1. Check spam folder
2. Verify FROM_EMAIL is a valid address
3. Check SMTP credentials
4. Review email logs

## Next Steps

- [Authentication Guide](authentication.md) - Password reset flow
- [Security Guide](../reference/security.md) - Email security best practices
- [Testing Guide](testing.md) - Testing email functionality
