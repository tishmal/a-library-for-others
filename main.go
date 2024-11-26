package main

import (
	"a-library-for-others/csvparser"
	"fmt"
	"io"
	"os"
)

func main() {
	// Открываем CSV файл
	file, err := os.Open("example.csv")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	var csvparser csvparser.CSVParser

	// Чтение файла построчно
	for {
		line, err := csvparser.ReadLine(file)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Error reading line:", err)
			return
		}

		// Печать строк и полей
		fmt.Println("Line:", line)
		for i := 0; i < csvparser.GetNumberOfFields(); i++ {
			field, err := csvparser.GetField(i)
			if err != nil {
				fmt.Println("Error getting field:", err)
				return
			}
			fmt.Printf("Field %d: %s\n", i, field)
		}
	}
}
