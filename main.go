package main

import (
	"errors"
	"fmt"
	"io"
	"os"
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
	lastLine   []byte   // последняя прочитанная строка
	lastFields [][]byte // поля в последней строке
	numFields  int      // количество полей в последней строке
}

// Вспомогательная функция для удаления \r и \n справа
func trimRight(line []byte) []byte {
	for len(line) > 0 && (line[len(line)-1] == '\r' || line[len(line)-1] == '\n') {
		line = line[:len(line)-1]
	}
	return line
}

// Вспомогательная функция проверки наличия символа в байтовом массиве
func contains(line []byte, char byte) bool {
	for _, b := range line {
		if b == char {
			return true
		}
	}
	return false
}

// Вспомогательная функция проверки начала и конца строки
func hasPrefix(line []byte, prefix byte) bool {
	return len(line) > 0 && line[0] == prefix
}

func hasSuffix(line []byte, suffix byte) bool {
	return len(line) > 0 && line[len(line)-1] == suffix
}

// Реализация метода ReadLine для структуры csvParser
func (c *csvParser) ReadLine(r io.Reader) (string, error) {
	var line []byte
	buf := make([]byte, 1)
	// Создаем объект, реализующий интерфейс CSVParser
	var csvparser CSVParser = &csvParser{}

	for {
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			return "", err
		}
		if n == 0 || err == io.EOF {
			// Если строка не пустая, обработаем последнюю строку
			if len(line) > 0 {
				result := trimRight(line)

				// Проверяем на лишние или отсутствующие кавычки
				if contains(result, '"') {
					if !hasPrefix(result, '"') || !hasSuffix(result, '"') {
						return "", ErrQuote
					}
				}

				// Обновляем состояние последней строки
				c.lastLine = result
				// Разделяем строку на поля
				c.lastFields = parseFields(result)
				c.numFields = len(c.lastFields)

				// Обновляем поля с помощью parseFields
				csvparser.(*csvParser).lastFields = parseFields([]byte(line))
				csvparser.(*csvParser).numFields = len(csvparser.(*csvParser).lastFields)

				// Выводим последнюю строку
				fmt.Println("Read line:", string(line))

				// Пример вывода количества полей
				numFields := csvparser.GetNumberOfFields()
				fmt.Printf("Number of fields: %d\n", numFields)

				// Пример вывода всех полей
				for i := 0; i < numFields; i++ {
					field, err := csvparser.GetField(i)
					if err != nil {
						fmt.Println("Error getting field:", err)
						os.Exit(1)
					}
					fmt.Printf("Field %d: %s\n", i, string(field))
				}

				return string(result), io.EOF
			}

			return "", io.EOF
		}

		// Добавляем символ в строку
		if buf[0] == '\n' {
			result := trimRight(line)

			// Проверяем на лишние или отсутствующие кавычки
			if contains(result, '"') {
				if !hasPrefix(result, '"') || !hasSuffix(result, '"') {
					return "", ErrQuote
				}
			}

			// Обновляем состояние последней строки
			c.lastLine = result
			// Разделяем строку на поля
			c.lastFields = parseFields(result)
			c.numFields = len(c.lastFields)

			return string(result), nil
		}
		line = append(line, buf[0])
	}
}

// parseFields разбивает строку на поля по запятым с учетом кавычек
func parseFields(line []byte) [][]byte {
	var fields [][]byte
	var currentField []byte
	inQuotes := false

	for i := 0; i < len(line); i++ {
		char := line[i]

		if char == '"' {
			inQuotes = !inQuotes // переключаем состояние кавычек
		} else if char == ',' && !inQuotes {
			// если мы не внутри кавычек, то это конец поля
			fields = append(fields, currentField)
			currentField = nil
		} else {
			// добавляем символ в текущее поле
			currentField = append(currentField, char)
		}
	}

	// Добавляем последнее поле
	fields = append(fields, currentField)

	return fields
}

// Реализация метода GetField для структуры csvParser
func (c *csvParser) GetField(n int) (string, error) {
	if n < 0 || n >= c.numFields {
		return "", ErrFieldCount
	}
	return string(c.lastFields[n]), nil
}

// Реализация метода GetNumberOfFields для структуры csvParser
func (c *csvParser) GetNumberOfFields() int {
	return c.numFields
}

func main() {
	args := os.Args
	if len(args) > 1 {
		fmt.Println("Incorrect input")
		os.Exit(1)
	}
	// Открываем файл
	file, err := os.Open("example.csv")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Создаем объект, реализующий интерфейс CSVParser
	var csvparser CSVParser = &csvParser{}

	// Чтение строк из CSV файла
	for {
		// Чтение строки с помощью ReadLine
		line, err := csvparser.ReadLine(file)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Error reading line:", err)
			return
		}

		// Выводим строку
		fmt.Println("Read line:", string(line))

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
			fmt.Printf("Field %d: %s\n", i, string(field))
		}
	}
}
