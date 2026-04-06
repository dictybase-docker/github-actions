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
	StockData      StockData
}

// StrainInfo represents information about a strain from the order.
type StrainInfo struct {
	ID         string
	Descriptor string
}

// PlasmidInfo represents information and storage details for a plasmid.
type PlasmidInfo struct {
	ID   string
	Name string
}

// StockData represents all stock-related information extracted from an order.
type StockData struct {
	StrainInfo  []StrainInfo
	PlasmidInfo []PlasmidInfo
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

// ExtractBillingEmail finds an h2 element with text "Ship To:" and extracts the email
// from a div element following it (at any nesting depth).
func ExtractBillingEmail(doc *html.Node) (string, error) {
	shipToH2 := findH2WithText(doc, "Ship To:")
	if shipToH2 == nil {
		return "", fmt.Errorf("'Ship To:' heading not found")
	}

	for sib := shipToH2.NextSibling; sib != nil; sib = sib.NextSibling {
		if email := findEmailInDiv(sib); email != "" {
			return email, nil
		}
	}

	return "", fmt.Errorf("billing address email not found")
}

// findH2WithText searches recursively for an h2 element with the given text content.
func findH2WithText(n *html.Node, text string) *html.Node {
	if n.Type == html.ElementNode && n.Data == "h2" {
		if strings.TrimSpace(getTextContent(n)) == text {
			return n
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if result := findH2WithText(c, text); result != nil {
			return result
		}
	}
	return nil
}

// findEmailInDiv recursively searches a node subtree for a div containing an email address.
func findEmailInDiv(n *html.Node) string {
	if n.Type == html.ElementNode && n.Data == "div" {
		if email := extractEmailFromText(getTextContent(n)); email != "" {
			return email
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if email := findEmailInDiv(c); email != "" {
			return email
		}
	}
	return ""
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

// containsIgnoreCase performs case-insensitive substring matching.
func containsIgnoreCase(str, substr string) bool {
	return strings.Contains(strings.ToLower(str), strings.ToLower(substr))
}

// matchesStrainInfoHeaders checks if headers match strain information table.
func matchesStrainInfoHeaders(headers []string) bool {
	if len(headers) < 5 {
		return false
	}

	hasID := false
	hasDescriptor := false
	hasStored := false

	for _, header := range headers {
		if containsIgnoreCase(header, "ID") {
			hasID = true
		}
		if containsIgnoreCase(header, "Descriptor") {
			hasDescriptor = true
		}
		if containsIgnoreCase(header, "Stored") {
			hasStored = true
		}
	}

	// Must have ID and Descriptor, but NOT Stored (to distinguish from plasmid storage)
	return hasID && hasDescriptor && !hasStored
}

// matchesPlasmidInfoHeaders checks if headers match plasmid information table.
func matchesPlasmidInfoHeaders(headers []string) bool {
	if len(headers) < 5 {
		return false
	}

	hasID := false
	hasName := false
	hasStored := false

	for _, header := range headers {
		if containsIgnoreCase(header, "ID") {
			hasID = true
		}
		if containsIgnoreCase(header, "Name") {
			hasName = true
		}
		if containsIgnoreCase(header, "Stored") {
			hasStored = true
		}
	}

	// Must have ID, Name, AND Stored (to distinguish from strain info which lacks Stored)
	return hasID && hasName && hasStored
}

// parseStrainInfoRow converts a table row to StrainInfo struct.
func parseStrainInfoRow(row []string) StrainInfo {
	strain := StrainInfo{}

	if len(row) > 0 {
		strain.ID = row[0]
	}
	if len(row) > 1 {
		strain.Descriptor = row[1]
	}

	return strain
}

// parsePlasmidInfoRow converts a table row to PlasmidInfo struct.
func parsePlasmidInfoRow(row []string) PlasmidInfo {
	plasmid := PlasmidInfo{}

	if len(row) > 0 {
		plasmid.ID = row[0]
	}
	if len(row) > 1 {
		plasmid.Name = row[1]
	}
	return plasmid
}

// isEmptyRow checks if a row contains only empty strings.
func isEmptyRow(row []string) bool {
	for _, cell := range row {
		if strings.TrimSpace(cell) != "" {
			return false
		}
	}
	return true
}

// ExtractStockData extracts strain and plasmid information from HTML content.
// It searches for tables with matching headers and parses stock-related data.
// Returns StockData with empty slices if no matching tables are found.
func ExtractStockData(doc *html.Node) (StockData, error) {
	tables, err := ParseTables(doc)
	if err != nil {
		return StockData{}, fmt.Errorf("error parsing tables: %w", err)
	}

	data := StockData{
		StrainInfo:  []StrainInfo{},
		PlasmidInfo: []PlasmidInfo{},
	}

	for _, table := range tables {
		// Check if this is a strain information table
		if matchesStrainInfoHeaders(table.Headers) {
			for _, row := range table.Rows {
				if !isEmptyRow(row) {
					data.StrainInfo = append(data.StrainInfo, parseStrainInfoRow(row))
				}
			}
		}

		// Check if this is a plasmid information table
		if matchesPlasmidInfoHeaders(table.Headers) {
			for _, row := range table.Rows {
				if !isEmptyRow(row) {
					data.PlasmidInfo = append(data.PlasmidInfo, parsePlasmidInfoRow(row))
				}
			}
		}
	}

	return data, nil
}

// ExtractOrderData extracts order ID, billing email, and stock data from HTML and returns structured data.
func ExtractOrderData(htmlNode *html.Node) (IssueBodyData, error) {
	orderID, err := ExtractOrderID(htmlNode)
	if err != nil {
		return IssueBodyData{}, fmt.Errorf("failed to extract order ID: %w", err)
	}

	billingEmail, err := ExtractBillingEmail(htmlNode)
	if err != nil {
		return IssueBodyData{}, fmt.Errorf("failed to extract billing email: %w", err)
	}

	stockData, err := ExtractStockData(htmlNode)
	if err != nil {
		return IssueBodyData{}, fmt.Errorf("failed to extract stock data: %w", err)
	}

	return IssueBodyData{
		OrderID:        orderID,
		RecipientEmail: billingEmail,
		StockData:      stockData,
	}, nil
}
