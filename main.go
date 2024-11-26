package main

import (
	"a-library-for-others/pkg/parser"
	"fmt"
	"io"
	"os"
)

func main() {
	file, err := os.Open("example.csv")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	var csvparser CSVParser = &CSVParserImpl{}

	for {
		line, err := csvparser.ReadLine(file)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Error reading line:", err)
			return
		}
		fmt.Println("Line:", line)

		numberOfFields := csvparser.GetNumberOfFields()
		fmt.Println("Number of fields:", numberOfFields)

		for i := 0; i < numberOfFields; i++ {
			field, err := csvparser.GetField(i)
			if err != nil {
				fmt.Println("Error getting field:", err)
				continue
			}
			fmt.Printf("Field %d: %s\n", i, field)
		}
	}
}
