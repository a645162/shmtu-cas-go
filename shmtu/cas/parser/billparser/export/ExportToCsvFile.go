package export

import (
	"encoding/csv"
	"fmt"
	"os"
	"shmtu-cas-go/shmtu/cas/parser/billparser"
)

func ToCsvFile(
	savePath string,
	billList []billparser.BillItemInfo,
) error {
	if billList == nil {
		return fmt.Errorf("billList is nil")
	}
	if len(billList) == 0 {
		return fmt.Errorf("billList is empty")
	}

	if savePath == "" {
		savePath = "./result.csv"
	}

	file, err := os.Create(savePath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	writer := csv.NewWriter(file)

	// Convert To Hashmap List
	data := make([]map[string]string, 0, len(billList))
	for _, bill := range billList {
		result := billparser.ConvertBillInfoToHashmap(&bill)
		data = append(data, result)
	}

	// Write headers
	headers := make([]string, 0, len(data[0]))
	for k := range data[0] {
		headers = append(headers, k)
	}
	err = writer.Write(headers)
	if err != nil {
		panic(err)
	}

	// Write data
	for _, row := range data {
		csvRow := make([]string, len(headers))
		for i, header := range headers {
			csvRow[i] = row[header]
		}
		err := writer.Write(csvRow)
		if err != nil {
			panic(err)
		}
	}

	// Call Flush to ensure all data written to the underlying writer
	writer.Flush()

	// Check for any error that occurred during the Flush.
	if err := writer.Error(); err != nil {
		return err
	}

	return nil
}
