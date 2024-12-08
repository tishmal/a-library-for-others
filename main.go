package main

import (
	"fmt"
	"io"
	"os"

	"a-library-for-others/csv"
)

func main() {
	// Открываем файл
	file, err := os.Open("example.csv")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	stat, err := file.Stat()
	if err != nil {
		fmt.Println("information: could not get file stats")
		os.Exit(1)
	}
	if stat.Size() == 0 {
		fmt.Println("information: file is empty")
		os.Exit(1)
	}
	defer file.Close()

	// Создаем объект, реализующий интерфейс CSVParser
	var csvparser csv.CSVParser = &csv.CsvParser{}

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
			fmt.Printf("Field %d: %s\n", i, string(field))
		}
		fmt.Println()
	}
}
