package csvparser

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

var (
	ErrQuote      = errors.New("excess or missing \" in quoted-field")
	ErrFieldCount = errors.New("wrong number of fields")
)

type CSVParser struct {
	lastLine   string
	fieldCount int
	lastFields []string
}

// ReadLine reads a single line from the input stream.
func (p *CSVParser) ReadLine(r io.Reader) (string, error) {
	scanner := bufio.NewScanner(r)
	if scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimRight(line, "\r\n")
		// If there's an unbalanced quote in the line, return an error.
		if strings.Count(line, "\"")%2 != 0 {
			return "", ErrQuote
		}
		p.lastLine = line
		p.lastFields = p.parseFields(line)
		p.fieldCount = len(p.lastFields)
		return line, nil
	}
	if scanner.Err() != nil {
		return "", scanner.Err()
	}
	return "", io.EOF
}

// parseFields splits a CSV line into fields, handling quotes.
func (p *CSVParser) parseFields(line string) []string {
	var fields []string
	var field strings.Builder
	inQuotes := false
	for i := 0; i < len(line); i++ {
		char := line[i]
		if char == '"' {
			if inQuotes && i+1 < len(line) && line[i+1] == '"' {
				field.WriteByte('"') // handle doubled quotes inside quoted fields
				i++                  // skip the next quote
			} else {
				inQuotes = !inQuotes
			}
		} else if char == ',' && !inQuotes {
			fields = append(fields, field.String())
			field.Reset()
		} else {
			field.WriteByte(char)
		}
	}
	fields = append(fields, field.String()) // append the last field
	return fields
}

// GetField returns the nth field of the last line read.
func (p *CSVParser) GetField(n int) (string, error) {
	if n < 0 || n >= p.fieldCount {
		return "", ErrFieldCount
	}
	return p.lastFields[n], nil
}

// GetNumberOfFields returns the number of fields in the last line read.
func (p *CSVParser) GetNumberOfFields() int {
	return p.fieldCount
}
