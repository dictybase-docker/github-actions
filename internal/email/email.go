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

// MailgunClient wraps the Mailgun client.
type MailgunClient struct {
	mg     *mailgun.MailgunImpl
	config MailgunConfig
}

// NewEmailClient creates a new email client with Mailgun configuration.
func NewEmailClient(domain, apiKey, from string) *MailgunClient {
	mg := mailgun.NewMailgun(domain, apiKey)
	return &MailgunClient{
		mg: mg,
		config: MailgunConfig{
			Domain: domain,
			APIKey: apiKey,
			From:   from,
		},
	}
}

// SendOrderUpdateEmail sends an order update email to the recipient.
func (ec *MailgunClient) SendOrderUpdateEmail(
	ctx context.Context,
	recipient string,
	subject string,
	htmlBody string,
) error {
	// Create a new message using the package function instead of method
	message := mailgun.NewMessage(
		ec.config.From,
		subject,
		"", // Plain text body (empty, using HTML only)
		recipient,
	)

	// Set HTML body
	message.SetHTML(htmlBody)

	// Send the message with a 10 second timeout
	sendCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	resp, messageID, err := ec.mg.Send(sendCtx, message)
	if err != nil {
		return fmt.Errorf("failed to send email via Mailgun: %w", err)
	}

	// Log success (could use logger here if needed)
	_ = resp      // Response body
	_ = messageID // Message ID

	return nil
}

// SendOrderUpdateFromTemplate sends an order update email using the template.
func (ec *MailgunClient) SendOrderUpdateFromTemplate(
	ctx context.Context,
	recipient string,
	data OrderEmailData,
) error {
	// Get the template path (relative to the email package)
	templatePath := filepath.Join("internal", "email", "order_update.tmpl")

	// Parse the template file
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", templatePath, err)
	}

	// Execute template with data
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// Create subject line
	subject := fmt.Sprintf("Dicty Stock Center - Order Update #%s", data.OrderID)

	// Send the email
	if err := ec.SendOrderUpdateEmail(ctx, recipient, subject, buf.String()); err != nil {
		return fmt.Errorf("failed to send order update email: %w", err)
	}

	return nil
}
