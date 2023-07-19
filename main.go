package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/PuerkitoBio/goquery"
	progressbar "github.com/cheggaaa/pb/v3"
)

func crawlDataByID(id string) ([]float64, error) {
	url := "https://vietnamnet.vn/giao-duc/diem-thi/tra-cuu-diem-thi-tot-nghiep-thpt/2023/" + id + ".html"

	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("URL returned 404")
	}

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, err
	}

	subjectMap := map[string]float64{
		"Toán":      0.0,
		"Văn":       0.0,
		"Ngoại ngữ": 0.0,
		"Lý":        0.0,
		"Hóa":       0.0,
		"Sinh":      0.0,
		"Sử":        0.0,
		"Địa":       0.0,
		"GDCD":      0.0,
	}

	keys := []string{"Toán", "Văn", "Ngoại ngữ", "Lý", "Hóa", "Sinh", "Sử", "Địa", "GDCD"}

	for _, key := range keys {
		valueStr := doc.Find("table tbody tr:contains('" + key + "') td:nth-child(2)").Text()
		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			continue
		}
		subjectMap[key] = value
	}

	values := make([]float64, len(subjectMap))
	for i, key := range keys {
		values[i] = subjectMap[key]
	}

	return values, nil
}

func insertData(id int, data []float64, writer *csv.Writer) error {
	stringData := make([]string, len(data)+1)
	stringData[0] = strconv.Itoa(id)

	for i, val := range data {
		stringData[i+1] = strconv.FormatFloat(val, 'f', -1, 64)
	}

	err := writer.Write(stringData)
	if err != nil {
		return fmt.Errorf("failed to write data: %w", err)
	}

	return nil
}

func updateData(file *os.File) {
	var wg sync.WaitGroup
	idCh := make(chan int)

	numThreads := 3000
	start := 1000000
	end := 640000000
	bar := progressbar.StartNew(end - start + 1)

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for i := 0; i < numThreads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for id := range idCh {
				var sid string
				if id < 10000000 {
					sid = "0" + strconv.Itoa(id)
				} else {
					sid = strconv.Itoa(id)
				}
				data, err := crawlDataByID(sid)
				if err != nil {
					continue
				}
				err = insertData(id, data, writer)
				if err != nil {
					fmt.Println(err)
				}
				bar.Increment()
			}
		}()
	}

	for id := start; id < end; id++ {
		idCh <- id
	}

	close(idCh)

	wg.Wait()
	bar.Finish()

	fmt.Println("Data inserted successfully!")
}

func main() {
	flag.Parse()

	filePath := "dataTHPT2023.csv"
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		fmt.Printf("failed to open file: %v", err)
		return
	}
	defer file.Close()

	updateData(file)

	fmt.Println("Data inserted successfully!")
}
