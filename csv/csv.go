package csv

import (
	"fmt"
	"io"
)

var (
	ErrQuote      = New("excess or missing \" in quoted-field")
	ErrFieldCount = New("wrong number of fields")
)

// Интерфейс CSVParser
type CSVParser interface {
	ReadLine(r io.Reader) (string, error)
	GetField(n int) (string, error)
	GetNumberOfFields() int
}

// Структура, реализующая интерфейс CSVParser
type CsvParser struct {
	lastLine   []byte   // последняя прочитанная строка
	lastFields [][]byte // поля в последней строке
	numFields  int      // количество полей в последней строке
}

// errors:
func New(text string) error {
	return &errorString{text}
}

type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

func isEmpty(r io.Reader) (bool, error) {
	buf := make([]byte, 1)
	_, err := r.Read(buf)
	if err == io.EOF {
		return true, nil
	}
	if err != nil {
		return false, err
	}
	return false, nil
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
func (c *CsvParser) ReadLine(r io.Reader) (string, error) {
	var line []byte
	buf := make([]byte, 1)
	// Создаем объект, реализующий интерфейс CSVParser
	var csvparser CSVParser = &CsvParser{}
	var count int
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
					for i := 0; i < len(line); i++ {
						count++
						if line[i] == '"' {
							line[i] = 0
						}
					}

					// if !hasPrefix(result, '"') || !hasSuffix(result, '"') {
					// 	return "", ErrQuote
					// }
				}
				if count%2 == 0 {
					fmt.Println(string(result))
				} else {
					return "", ErrQuote
				}

				// Обновляем состояние последней строки
				c.lastLine = result
				// Разделяем строку на поля
				c.lastFields = parseFields(result)
				c.numFields = len(c.lastFields)

				// Обновляем поля с помощью parseFields
				csvparser.(*CsvParser).lastFields = parseFields([]byte(line))
				csvparser.(*CsvParser).numFields = len(csvparser.(*CsvParser).lastFields)

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
				for i := 0; i < len(line); i++ {
					count++
					if line[i] == '"' {
						line[i] = 0
					}
				}

				// if !hasPrefix(result, '"') || !hasSuffix(result, '"') {
				// 	return "", ErrQuote
				// }
			}
			if count%2 == 0 {
				fmt.Println(string(result))
			} else {
				return "", ErrQuote
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
func (c *CsvParser) GetField(n int) (string, error) {
	if n < 0 || n >= c.numFields {
		return "", ErrFieldCount
	}
	return string(c.lastFields[n]), nil
}

// Реализация метода GetNumberOfFields для структуры csvParser
func (c *CsvParser) GetNumberOfFields() int {
	return c.numFields
}
