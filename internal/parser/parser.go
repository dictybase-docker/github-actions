package parser

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"golang.org/x/net/html"
)

// TableData represents a parsed HTML table.
type TableData struct {
	Headers []string
	Rows    [][]string
}

type IssueBodyData struct {
	RecipientEmail string
	OrderID        string
}

var (
	emailRe   = regexp.MustCompile(`([a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,})`)
	orderIDRe = regexp.MustCompile(`Order\s*ID:\s*(\d+)`)
)

// ParseTables extracts all tables from HTML content.
func ParseTables(doc *html.Node) ([]TableData, error) {
	var tables []TableData
	var findTables func(*html.Node)
	findTables = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "table" {
			tables = append(tables, parseTable(n))
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findTables(c)
		}
	}
	findTables(doc)

	return tables, nil
}

// parseTable extracts data from a single table node.
func parseTable(tableNode *html.Node) TableData {
	var table TableData

	for child := tableNode.FirstChild; child != nil; child = child.NextSibling {
		if child.Type != html.ElementNode {
			continue
		}

		switch child.Data {
		case "thead":
			table.Headers = parseTableHead(child)
		case "tbody":
			table.Rows = parseTableBody(child)
		case "tr":
			// Table without thead/tbody - first row becomes headers
			cells := parseTableRow(child)
			if len(table.Headers) == 0 {
				table.Headers = cells
			} else {
				table.Rows = append(table.Rows, cells)
			}
		}
	}

	return table
}

// parseTableHead extracts header cells from thead.
func parseTableHead(theadNode *html.Node) []string {
	for row := theadNode.FirstChild; row != nil; row = row.NextSibling {
		if row.Type == html.ElementNode && row.Data == "tr" {
			return parseTableRow(row)
		}
	}
	return nil
}

// parseTableBody extracts all rows from tbody.
func parseTableBody(tbodyNode *html.Node) [][]string {
	var rows [][]string
	for row := tbodyNode.FirstChild; row != nil; row = row.NextSibling {
		if row.Type == html.ElementNode && row.Data == "tr" {
			rows = append(rows, parseTableRow(row))
		}
	}
	return rows
}

// parseTableRow extracts cells from a table row.
func parseTableRow(rowNode *html.Node) []string {
	var cells []string
	for cell := rowNode.FirstChild; cell != nil; cell = cell.NextSibling {
		if cell.Type == html.ElementNode && (cell.Data == "td" || cell.Data == "th") {
			cells = append(cells, getTextContent(cell))
		}
	}
	return cells
}

// getTextContent recursively extracts all text from a node.
func getTextContent(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	var text strings.Builder
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		_, _ = text.WriteString(getTextContent(c))
	}
	return strings.TrimSpace(text.String())
}

// ExtractBillingEmail finds a table with "Billing Address" header and extracts the email from that column.
func ExtractBillingEmail(doc *html.Node) (string, error) {
	tables, err := ParseTables(doc)
	if err != nil {
		return "", fmt.Errorf("error parsing tables: %w", err)
	}

	// Find table with "Billing Address" header
	for _, table := range tables {
		billingColIndex := -1

		// Search for "Billing Address" header (case-insensitive)
		for i, header := range table.Headers {
			if strings.Contains(strings.ToLower(header), "billing") &&
				strings.Contains(strings.ToLower(header), "address") {
				billingColIndex = i
				break
			}
		}

		// If this table has the billing address column
		if billingColIndex != -1 {
			// Search all rows for an email in the billing address column
			for _, row := range table.Rows {
				if billingColIndex < len(row) {
					email := extractEmailFromText(row[billingColIndex])
					if email != "" {
						return email, nil
					}
				}
			}
		}
	}

	return "", fmt.Errorf("billing address email not found")
}

// extractEmailFromText extracts an email address from a string using regex.
func extractEmailFromText(text string) string {
	matches := emailRe.FindStringSubmatch(text)
	if len(matches) >= 2 {
		return matches[1]
	}

	return ""
}

// ExtractOrderID finds a paragraph containing "Order ID" and extracts the order ID value.
func ExtractOrderID(doc *html.Node) (string, error) {
	var orderID string
	var findOrderParagraph func(*html.Node)
	findOrderParagraph = func(node *html.Node) {
		if orderID != "" {
			return // Already found
		}

		if node.Type == html.ElementNode && node.Data == "p" {
			text := getTextContent(node)
			extracted := extractOrderIDFromText(text)
			if extracted != "" {
				orderID = extracted
			}
		}

		for c := node.FirstChild; c != nil; c = c.NextSibling {
			findOrderParagraph(c)
		}
	}
	findOrderParagraph(doc)

	if orderID == "" {
		return "", fmt.Errorf("order ID not found in paragraphs")
	}

	return orderID, nil
}

// extractOrderIDFromText extracts an order ID from text using regex.
func extractOrderIDFromText(text string) string {
	matches := orderIDRe.FindStringSubmatch(text)
	if len(matches) >= 2 {
		return matches[1]
	}

	return ""
}

// MarkdownToHTML converts markdown text to HTML using goldmark with GitHub Flavored Markdown extensions.
func MarkdownToHTML(markdown string) (*html.Node, error) {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
	)

	var buf bytes.Buffer
	if err := md.Convert([]byte(markdown), &buf); err != nil {
		return nil, fmt.Errorf("error converting markdown: %w", err)
	}

	doc, err := html.Parse(strings.NewReader(buf.String()))
	if err != nil {
		return nil, fmt.Errorf("error parsing HTML: %w", err)
	}

	return doc, nil
}

// ExtractOrderData extracts order ID and billing email from HTML and returns structured data.
func ExtractOrderData(htmlNode *html.Node) (IssueBodyData, error) {
	orderID, err := ExtractOrderID(htmlNode)
	if err != nil {
		return IssueBodyData{}, fmt.Errorf("failed to extract order ID: %w", err)
	}

	billingEmail, err := ExtractBillingEmail(htmlNode)
	if err != nil {
		return IssueBodyData{}, fmt.Errorf("failed to extract billing email: %w", err)
	}

	return IssueBodyData{
		OrderID:        orderID,
		RecipientEmail: billingEmail,
	}, nil
}
