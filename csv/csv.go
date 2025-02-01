package csv

import (
	"encoding/csv"
	"fmt"
	"os"
)

type PhoneColumns struct {
	Headers     []string
	PhoneFields map[string]int
}

type CSVProcessor struct {
	reader      *csv.Reader
	writer      *csv.Writer
	inputFile   *os.File
	outputFile  *os.File
	phoneFields PhoneColumns
}

func NewCSVProcessor(inputPath, outputPath string) (*CSVProcessor, error) {
	inFile, err := os.Open(inputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open input file: %w", err)
	}

	outFile, err := os.Create(outputPath)
	if err != nil {
		inFile.Close()
		return nil, fmt.Errorf("failed to create output file: %w", err)
	}

	processor := &CSVProcessor{
		reader:     csv.NewReader(inFile),
		writer:     csv.NewWriter(outFile),
		inputFile:  inFile,
		outputFile: outFile,
	}

	if err := processor.initializeHeaders(); err != nil {
		processor.Close()
		return nil, err
	}

	return processor, nil
}

func (p *CSVProcessor) initializeHeaders() error {
	headers, err := p.reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read headers: %w", err)
	}

	p.phoneFields = PhoneColumns{
		Headers:     headers,
		PhoneFields: make(map[string]int),
	}

	for i, header := range headers {
		switch header {
		case "Primary #", "Phone 1", "Phone 2", "Phone 3":
			p.phoneFields.PhoneFields[header] = i
		}
	}

	// Write new headers with result columns
	newHeaders := headers
	for colName := range p.phoneFields.PhoneFields {
		newHeaders = append(newHeaders, colName+" Result")
	}
	return p.writer.Write(newHeaders)
}

func (p *CSVProcessor) ProcessRows(processor func(string) (bool, error)) error {
	for {
		record, err := p.reader.Read()
		if err != nil {
			break
		}

		newRow := record
		for _, idx := range p.phoneFields.PhoneFields {
			phoneNum := record[idx]
			result, err := processor(phoneNum)
			if err != nil {
				newRow = append(newRow, "Error")
				continue
			}
			if result {
				newRow = append(newRow, "Found")
			} else {
				newRow = append(newRow, "Not Found")
			}
		}
		if err := p.writer.Write(newRow); err != nil {
			return err
		}
	}
	return nil
}

func (p *CSVProcessor) Close() error {
	p.writer.Flush()
	p.inputFile.Close()
	return p.outputFile.Close()
}
