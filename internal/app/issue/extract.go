package issue

import (
	"fmt"
	"regexp"
)

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
func extractOrderID(text string) (string, error) {
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

func extractEmail(text string) (string, error) {
	// Standard email pattern
	pattern := `([a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,})`
	re := regexp.MustCompile(pattern)

	matches := re.FindStringSubmatch(text)
	if len(matches) < 2 {
		return "", fmt.Errorf("email not found")
	}

	return matches[1], nil
}

func extractWithRegex(text string) (orderID string, email string, err error) {
	orderID, err = extractOrderID(text)
	if err != nil {
		return "", "", fmt.Errorf("failed to extract order ID: %w", err)
	}

	email, err = extractEmail(text)
	if err != nil {
		return "", "", fmt.Errorf("failed to extract email: %w", err)
	}

	return orderID, email, nil
}
