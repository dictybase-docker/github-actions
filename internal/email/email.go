package email

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"html/template"
	"time"

	"github.com/mailgun/mailgun-go/v5"
)

//go:embed order_update.tmpl
var templateFS embed.FS

// OrderEmailData represents the data structure for order update emails.
type OrderEmailData struct {
	RecipientEmail string
	OrderID        string
	Label          string
}

// MailgunConfig holds Mailgun configuration.
type MailgunConfig struct {
	Domain string
	From   string // Sender email address
}

// MailgunClient wraps the Mailgun client.
type MailgunClient struct {
	mg     *mailgun.Client
	config MailgunConfig
}

// NewEmailClient creates a new email client with Mailgun configuration.
func NewEmailClient(domain, apiKey, from string) *MailgunClient {
	mg := mailgun.NewMailgun(apiKey)
	return &MailgunClient{
		mg: mg,
		config: MailgunConfig{
			Domain: domain,
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
	// Create a new message - domain is now passed to NewMessage in v5
	message := mailgun.NewMessage(
		ec.config.Domain,
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

	_, err := ec.mg.Send(sendCtx, message)
	if err != nil {
		return fmt.Errorf("failed to send email via Mailgun: %w", err)
	}

	return nil
}

func createEmailHTML(data OrderEmailData) (string, error) {
	// Parse the embedded template file
	tmpl, err := template.ParseFS(templateFS, "order_update.tmpl")
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	// Execute template with data
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}
	return buf.String(), nil
}

// SendOrderUpdateFromTemplate sends an order update email using the template.
func (ec *MailgunClient) SendOrderUpdateFromTemplate(
	ctx context.Context,
	recipient string,
	data OrderEmailData,
) error {
	html, err := createEmailHTML(data)

	if err != nil {
		return fmt.Errorf("failed to create email HTML: %w", err)
	}

	// Create subject line
	subject := fmt.Sprintf("Dicty Stock Center - Order Update #%s", data.OrderID)

	// Send the email
	if err := ec.SendOrderUpdateEmail(ctx, recipient, subject, html); err != nil {
		return fmt.Errorf("failed to send order update email: %w", err)
	}

	return nil
}
