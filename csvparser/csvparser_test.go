package csvparser

import (
	"strings"
	"testing"
)

func TestCSVParser(t *testing.T) {
	csv := &CSVParser{}

	// Тест для строки без кавычек
	line := "a,b,c"
	r := strings.NewReader(line)
	_, err := csv.ReadLine(r)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Проверка правильности количества полей
	if csv.GetNumberOfFields() != 3 {
		t.Fatalf("expected 3 fields, got %d", csv.GetNumberOfFields())
	}

	// Тест для строки с кавычками
	line = `"a,b",c,"d,e"`
	r = strings.NewReader(line)
	_, err = csv.ReadLine(r)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if csv.GetNumberOfFields() != 3 {
		t.Fatalf("expected 3 fields, got %d", csv.GetNumberOfFields())
	}

	// Тест на ошибку с неправильными кавычками
	line = `"a,b,c`
	r = strings.NewReader(line)
	_, err = csv.ReadLine(r)
	if err != ErrQuote {
		t.Fatalf("expected ErrQuote, got %v", err)
	}
}
