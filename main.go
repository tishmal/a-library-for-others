package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	ErrQuote      = errors.New("excess or missing \" in quoted-field")
	ErrFieldCount = errors.New("wrong number of fields")
)

// Интерфейс CSVParser
type CSVParser interface {
	ReadLine(r io.Reader) (string, error)
	GetField(n int) (string, error)
	GetNumberOfFields() int
}

// Структура, реализующая интерфейс CSVParser
type csvParser struct {
	lastLine   string   // последняя прочитанная строка
	lastFields []string // поля в последней строке
	numFields  int      // количество полей в последней строке
}

// Реализация метода ReadLine для структуры csvParser
func (c *csvParser) ReadLine(r io.Reader) (string, error) {
	// Чтение строки
	bufReader := bufio.NewReader(r)

	// Читаем строку
	line, err := bufReader.ReadString('\n')
	if err != nil {
		// Если ошибка EOF, возвращаем пустую строку и nil (чтобы завершить чтение)
		if err == io.EOF {
			return "", io.EOF
		}
		return "", err
	}

	// Убираем символ новой строки
	line = strings.TrimRight(line, "\r\n")

	// Проверяем на лишние или отсутствующие кавычки
	if strings.Contains(line, "\"") {
		if !strings.HasPrefix(line, "\"") || !strings.HasSuffix(line, "\"") {
			return "", ErrQuote
		}
	}

	// Обновляем состояние последней строки
	c.lastLine = line
	// Разделяем строку на поля
	c.lastFields = parseFields(line)
	c.numFields = len(c.lastFields)

	// Возвращаем строку
	return line, nil
}

// parseFields разбивает строку на поля по запятым с учетом кавычек
func parseFields(line string) []string {
	var fields []string
	var currentField strings.Builder
	inQuotes := false

	for i := 0; i < len(line); i++ {
		char := line[i]

		if char == '"' {
			inQuotes = !inQuotes // переключаем состояние кавычек
		} else if char == ',' && !inQuotes {
			// если мы не внутри кавычек, то это конец поля
			fields = append(fields, currentField.String())
			currentField.Reset()
		} else {
			// добавляем символ в текущее поле
			currentField.WriteByte(char)
		}
	}

	// Добавляем последнее поле
	fields = append(fields, currentField.String())

	return fields
}

// Реализация метода GetField для структуры csvParser
func (c *csvParser) GetField(n int) (string, error) {
	if n < 0 || n >= c.numFields {
		return "", ErrFieldCount
	}
	return c.lastFields[n], nil
}

// Реализация метода GetNumberOfFields для структуры csvParser
func (c *csvParser) GetNumberOfFields() int {
	return c.numFields
}

func main() {
	// Открываем файл
	file, err := os.Open("example.csv")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Создаем объект, реализующий интерфейс CSVParser
	var csvparser CSVParser = &csvParser{}

	// Создаем буферизованный ридер для всего файла
	bufReader := bufio.NewReader(file)

	// Чтение строк из CSV файла
	for {
		// Чтение строки с помощью ReadLine
		line, err := csvparser.ReadLine(bufReader)
		if err != nil {
			if err == io.EOF {
				break // Конец файла, выходим
			}
			fmt.Println("Error reading line:", err)
			return
		}

		// Выводим строку
		fmt.Println("Read line:", line)

		// Пример вывода количества полей
		numFields := csvparser.GetNumberOfFields()
		fmt.Printf("Number of fields: %d\n", numFields)

		// Пример вывода всех полей
		for i := 0; i < numFields; i++ {
			field, err := csvparser.GetField(i)
			if err != nil {
				fmt.Println("Error getting field:", err)
				return
			}
			fmt.Printf("Field %d: %s\n", i, field)
		}
	}
}
