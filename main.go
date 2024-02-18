package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"regexp"
)

// Category represents a category for bank transactions
type Category struct {
	Name  string
	Regex *regexp.Regexp
}

func main() {
	if len(os.Args) != 5 {
		fmt.Println("Usage: go run main.go categories.csv transactions.csv output.csv")
		return
	}

	categoriesFile, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println("Error opening categories file:", err)
		return
	}
	defer categoriesFile.Close()

	categoriesReader := csv.NewReader(categoriesFile)
	categoriesRecords, err := categoriesReader.ReadAll()
	if err != nil {
		fmt.Println("Error reading categories file:", err)
		return
	}

	var categories []Category
	for _, record := range categoriesRecords {
		if len(record) == 2 {
			regex, err := regexp.Compile(record[1])
			if err != nil {
				fmt.Println("Error compiling regex for category", record[0], ":", err)
				continue
			}
			categories = append(categories, Category{Name: record[0], Regex: regex})
		}
	}

	transactionsFile, err := os.Open(os.Args[2])
	if err != nil {
		fmt.Println("Error opening transactions file:", err)
		return
	}
	defer transactionsFile.Close()

	transactionsReader := csv.NewReader(transactionsFile)
	transactionsRecords, err := transactionsReader.ReadAll()
	if err != nil {
		fmt.Println("Error reading transactions file:", err)
		return
	}

	var updatedTransactions [][]string
	for _, record := range transactionsRecords {
		var matchedCategory string
		for _, category := range categories {
			if category.Regex.MatchString(record[1]) {
				matchedCategory = category.Name
				break
			}
		}
		record = append(record, matchedCategory)
		updatedTransactions = append(updatedTransactions, record)
	}

	outputFile, err := os.Create(os.Args[3])
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer outputFile.Close()

	outputWriter := csv.NewWriter(outputFile)
	defer outputWriter.Flush()

	for _, record := range updatedTransactions {
		if err := outputWriter.Write(record); err != nil {
			fmt.Println("Error writing transaction data:", err)
			return
		}
	}

	fmt.Printf("Categories successfully added and transactions saved to '%s'.\n", os.Args[3])
}
