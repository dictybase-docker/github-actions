package issue

import (
	"fmt"
	"regexp"
)

type OrderData struct {
	orderID        string
	recipientEmail string
}

type IssueProcessor struct {
	issueBody string
	orderData OrderData
}

// extractOrderData parses the issueBody and populates the orderData field
func (ip *IssueProcessor) extractOrderData() error {
	// Use extraction logic to parse issueBody
	orderID, email, err := extractFromBody(ip.issueBody)
	if err != nil {
		return fmt.Errorf("failed to extract order data: %w", err)
	}

	// Write to the orderData field
	ip.orderData = OrderData{
		orderID:        orderID,
		recipientEmail: email,
	}

	return nil
}

func extractFromTemplate(text, templatePattern string) (map[string]string, error) {
	// Replace {{.Field}} with named capture groups
	regex := regexp.MustCompile(`\{\{\.(\w+)\}\}`)
	pattern := regex.ReplaceAllString(templatePattern, `(?P<$1>.+?)`)

	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(text)

	result := make(map[string]string)
	for i, name := range re.SubexpNames() {
		if i > 0 && i <= len(matches) {
			result[name] = matches[i]
		}
	}
	return result, nil
}

// ExtractOrderID extracts the Order ID from markdown text
// Pattern: Order ID: VALUE (with optional ** markdown bold)
func extractOrderIDFromTitle(text string) (string, error) {
	// Match: Order ID: followed by optional whitespace and the ID value
	// Handles both "**Order ID:**" and "Order ID:"
	pattern := `\*?\*?Order ID:\*?\*?\s*(\S+)`
	re := regexp.MustCompile(pattern)

	matches := re.FindStringSubmatch(text)
	if len(matches) < 2 {
		return "", fmt.Errorf("order ID not found")
	}

	return matches[1], nil
}

func extractEmailFromTitle(text string) (string, error) {
	// Standard email pattern
	pattern := `([a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,})`
	re := regexp.MustCompile(pattern)

	matches := re.FindStringSubmatch(text)
	if len(matches) < 2 {
		return "", fmt.Errorf("email not found")
	}

	return matches[1], nil
}

func extractFromTitle(text string) (orderID string, email string, err error) {
	orderID, err = extractOrderIDFromTitle(text)
	if err != nil {
		return "", "", fmt.Errorf("failed to extract order ID: %w", err)
	}

	email, err = extractEmailFromTitle(text)
	if err != nil {
		return "", "", fmt.Errorf("failed to extract email: %w", err)
	}

	return orderID, email, nil
}

// extractOrderIDFromBody extracts the Order ID from the issue body
// Pattern: **Order ID:** 37500885
func extractOrderIDFromBody(text string) (string, error) {
	pattern := `\*\*Order ID:\*\*\s+(\d+)`
	re := regexp.MustCompile(pattern)

	matches := re.FindStringSubmatch(text)
	if len(matches) < 2 {
		return "", fmt.Errorf("order ID not found in body")
	}

	return matches[1], nil
}

// extractBillingEmailFromBody extracts the email from the Billing address column
// The table structure is: | Shipping address | (empty) | Billing address |
// We find all emails and return the second one (billing email)
func extractBillingEmailFromBody(text string) (string, error) {
	// Find the billing section by looking for text after "Billing address"
	// Then extract email from that section
	billingIdx := regexp.MustCompile(`Billing address`).FindStringIndex(text)
	if billingIdx == nil {
		return "", fmt.Errorf("billing address section not found")
	}

	// Get text starting from billing address header
	billingSection := text[billingIdx[0]:]

	// Find the table row after the header (skip the separator line)
	// Look for a line with | ... | ... | content |
	lines := regexp.MustCompile(`\n`).Split(billingSection, -1)

	// Skip first two lines (header and separator), get the data row
	if len(lines) < 3 {
		return "", fmt.Errorf("billing data row not found")
	}

	dataRow := lines[2]

	// Split by | and get the third column (index 2)
	columns := regexp.MustCompile(`\|`).Split(dataRow, -1)
	if len(columns) < 4 {
		return "", fmt.Errorf("insufficient columns in billing row")
	}

	billingColumn := columns[3] // Third column is at index 3 (0=empty, 1=shipping, 2=middle, 3=billing)

	// Extract email from billing column
	emailPattern := `([a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,})`
	emailRe := regexp.MustCompile(emailPattern)

	matches := emailRe.FindStringSubmatch(billingColumn)
	if len(matches) < 2 {
		return "", fmt.Errorf("email not found in billing column")
	}

	return matches[1], nil
}

// extractFromBody extracts both Order ID and billing email from the issue body
func extractFromBody(text string) (orderID string, email string, err error) {
	orderID, err = extractOrderIDFromBody(text)
	if err != nil {
		return "", "", fmt.Errorf("failed to extract order ID: %w", err)
	}

	email, err = extractBillingEmailFromBody(text)
	if err != nil {
		return "", "", fmt.Errorf("failed to extract billing email: %w", err)
	}

	return orderID, email, nil
}
