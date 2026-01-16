package html

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

// TableData represents a parsed HTML table.
type TableData struct {
	Headers []string
	Rows    [][]string
}

// ParseTables extracts all tables from HTML content.
func ParseTables(htmlContent string) ([]TableData, error) {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return nil, fmt.Errorf("error parsing HTML: %w", err)
	}

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

// parseTable extracts data from a single table node
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

// parseTableHead extracts header cells from thead
func parseTableHead(theadNode *html.Node) []string {
	for row := theadNode.FirstChild; row != nil; row = row.NextSibling {
		if row.Type == html.ElementNode && row.Data == "tr" {
			return parseTableRow(row)
		}
	}
	return nil
}

// parseTableBody extracts all rows from tbody
func parseTableBody(tbodyNode *html.Node) [][]string {
	var rows [][]string
	for row := tbodyNode.FirstChild; row != nil; row = row.NextSibling {
		if row.Type == html.ElementNode && row.Data == "tr" {
			rows = append(rows, parseTableRow(row))
		}
	}
	return rows
}

// parseTableRow extracts cells from a table row
func parseTableRow(rowNode *html.Node) []string {
	var cells []string
	for cell := rowNode.FirstChild; cell != nil; cell = cell.NextSibling {
		if cell.Type == html.ElementNode && (cell.Data == "td" || cell.Data == "th") {
			cells = append(cells, getTextContent(cell))
		}
	}
	return cells
}

// getTextContent recursively extracts all text from a node
func getTextContent(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	var text strings.Builder
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		text.WriteString(getTextContent(c))
	}
	return strings.TrimSpace(text.String())
}
