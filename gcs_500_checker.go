package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	fptr := flag.String("filepath", "maybe500s.csv", "Absolute path to CSV of potential 500s exported from Google Search Console")
	
	flag.Parse()

	file, err := os.Open(*fptr)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	reader := csv.NewReader(file)

	rows, err := reader.ReadAll()

	if err != nil {
		log.Fatal(err)
	}

	newFile, err := os.Create("actual500s.csv")

	if err != nil {
		log.Fatal(err)
	}

	defer newFile.Close()

	writer := csv.NewWriter(newFile)

	defer writer.Flush()

	headers := []string{"url", "status"}

	writer.Write(headers)

	count := 0

	for idx, row := range rows {
		// ignore header
		if idx == 0 {
			continue
		}

		url := row[0]

		res, err := http.Get(url)

		if err != nil {
			fmt.Println(err)
			continue
		}

		defer res.Body.Close()

		statusCode := res.StatusCode

		if statusCode == 500 {
			newRow := []string{url, "500"}
			writer.Write(newRow)
			count++
		}
	}

	fmt.Printf("Done! Found %d actual 500 errors", count)
}