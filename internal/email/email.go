package email

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"path/filepath"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

// OrderEmailData represents the data structure for order update emails.
type OrderEmailData struct {
	RecipientEmail string
	OrderID        string
	Label          string
}

// MailgunConfig holds Mailgun configuration.
type MailgunConfig struct {
	Domain string
	APIKey string
	From   string // Sender email address
}

// EmailClient wraps the Mailgun client.
type EmailClient struct {
	mg     *mailgun.MailgunImpl
	config MailgunConfig
}

// NewEmailClient creates a new email client with Mailgun configuration
func NewEmailClient(domain, apiKey, from string) *EmailClient {
	mg := mailgun.NewMailgun(domain, apiKey)
	return &EmailClient{
		mg: mg,
		config: MailgunConfig{
			Domain: domain,
			APIKey: apiKey,
			From:   from,
		},
	}
}

// RenderTemplate renders the order update email template with provided data
func RenderTemplate(templatePath string, data OrderEmailData) (string, error) {
	// Parse the template file
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to parse template %s: %w", templatePath, err)
	}

	// Execute template with data
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// SendOrderUpdateEmail sends an order update email to the recipient
func (ec *EmailClient) SendOrderUpdateEmail(
	ctx context.Context,
	recipient string,
	subject string,
	htmlBody string,
) error {
	// Create a new message
	message := ec.mg.NewMessage(
		ec.config.From,
		subject,
		"", // Plain text body (empty, using HTML only)
		recipient,
	)

	// Set HTML body
	message.SetHtml(htmlBody)

	// Send the message with a 10 second timeout
	sendCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	resp, id, err := ec.mg.Send(sendCtx, message)
	if err != nil {
		return fmt.Errorf("failed to send email via Mailgun: %w", err)
	}

	// Log success (could use logger here if needed)
	_ = resp // Response body
	_ = id   // Message ID

	return nil
}

// SendOrderUpdateFromTemplate sends an order update email using the template
func (ec *EmailClient) SendOrderUpdateFromTemplate(
	ctx context.Context,
	recipient string,
	data OrderEmailData,
) error {
	// Get the template path (relative to the email package)
	templatePath := filepath.Join("internal", "email", "order_update.tmpl")

	// Render the template
	htmlBody, err := RenderTemplate(templatePath, data)
	if err != nil {
		return fmt.Errorf("failed to render email template: %w", err)
	}

	// Create subject line
	subject := fmt.Sprintf("Dicty Stock Center - Order Update #%s", data.OrderID)

	// Send the email
	if err := ec.SendOrderUpdateEmail(ctx, recipient, subject, htmlBody); err != nil {
		return fmt.Errorf("failed to send order update email: %w", err)
	}

	return nil
}
