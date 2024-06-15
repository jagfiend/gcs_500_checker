package main

import (
	"encoding/csv"
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	// open given file
	fptr := flag.String("filepath", "maybe500s.csv", "filepath of export")
	
	flag.Parse()

	file, err := os.Open(*fptr)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	// setup reader of given file
	reader := csv.NewReader(file)

	rows, err := reader.ReadAll()

	if err != nil {
		log.Fatal(err)
	}

	// setup new files
	actual500s, err := os.Create("actual500s.csv")

	if err != nil {
		log.Fatal(err)
	}

	defer actual500s.Close()

	not500s, err := os.Create("not500s.csv")

	if err != nil {
		log.Fatal(err)
	}

	defer not500s.Close()

	// setup writers
	writerActuals := csv.NewWriter(actual500s)

	defer writerActuals.Flush()

	writerNots := csv.NewWriter(not500s)

	defer writerNots.Flush()

	headers := []string{"url", "status"}

	writerActuals.Write(headers)
	writerNots.Write(headers)

	count500s := 0
	countNots := 0

	for idx, row := range rows {
		// ignore header
		if idx == 0 {
			continue
		}

		url := row[0]

		res, err := http.Get(url)

		if err != nil {
			log.Println(err)
			continue
		}

		defer res.Body.Close()

		statusCode := strconv.Itoa(res.StatusCode)
		newRow := []string{url, statusCode}

		if statusCode == "500" {
			writerActuals.Write(newRow)
			count500s++
		} else {
			writerNots.Write(newRow)
			countNots++
		}
	}

	log.Printf("Done! Found %d actual 500 errors and %d not", count500s, countNots)
}