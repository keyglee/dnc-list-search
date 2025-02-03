package csv

import (
	"dnclistsearch/logger"
	"encoding/csv"
	"fmt"
	"os"
	"slices"

	"github.com/sirupsen/logrus"
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
	logger      *logrus.Logger
}

func (c *CSVProcessor) GetReader() *csv.Reader {
	return c.reader
}

func NewCSVProcessor(inputPath, outputPath string, wantedHeaders []string) (*CSVProcessor, error) {
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
		logger:     logger.GetLogger(),
	}

	if err := processor.initializeHeaders(wantedHeaders); err != nil {
		processor.Close()
		return nil, err
	}

	return processor, nil
}

func (p *CSVProcessor) initializeHeaders(wantedHeaders []string) error {
	headers, err := p.reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read headers: %w", err)
	}

	p.phoneFields = PhoneColumns{
		Headers:     headers,
		PhoneFields: make(map[string]int),
	}

	for i, header := range headers {
		if slices.Contains(wantedHeaders, header) {
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

func (p *CSVProcessor) Close() error {
	p.writer.Flush()
	p.inputFile.Close()
	return p.outputFile.Close()
}
