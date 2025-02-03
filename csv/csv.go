package csv

import (
	"dnclistsearch/logger"
	"encoding/csv"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/sirupsen/logrus"
)

type PhoneColumns struct {
	Headers           []string
	WriteHeaders      []string
	SearchPhoneFields map[string]int
	WritePhoneFields  map[string]int
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

	processor.logger.Debug("CSVProcessor initialized")
	processor.logger.Debugf("Headers: %+v", processor.phoneFields.Headers)
	processor.logger.Debugf("Write Headers: %+v", processor.phoneFields.WriteHeaders)
	processor.logger.Debugf("SearchPhoneFields: %+v", processor.phoneFields.SearchPhoneFields)

	return processor, nil
}

func (p *CSVProcessor) initializeHeaders(wantedHeaders []string) error {
	headers, err := p.reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read headers: %w", err)
	}

	p.phoneFields = PhoneColumns{
		Headers:           headers,
		SearchPhoneFields: make(map[string]int),
	}

	for i, header := range headers {
		if slices.Contains(wantedHeaders, header) {
			p.phoneFields.SearchPhoneFields[header] = i
		}
	}

	// Write new headers with result columns
	newHeaders := headers
	for colName := range p.phoneFields.SearchPhoneFields {
		newHeaders = append(newHeaders, colName+" Result")

	}

	p.phoneFields.WriteHeaders = newHeaders

	return p.writer.Write(newHeaders)
}

func (p *CSVProcessor) GetPhoneNumbers(row []string) []string {
	phones := make([]string, 0)

	// Create a slice of keys and sort them
	keys := make([]int, 0, len(p.phoneFields.SearchPhoneFields))
	for _, col := range p.phoneFields.SearchPhoneFields {
		keys = append(keys, col)
	}
	slices.Sort(keys)

	// Append phone numbers in sorted order
	for _, col := range keys {
		phones = append(phones, row[col])
	}

	p.logger.Debugf("Extracted phone numbers: %v", phones)
	p.logger.Debugf("Number of phone numbers: %d", len(phones))

	return phones
}

func (p *CSVProcessor) WriteRow(row []string) error {

	p.logger.Debug("Writing row")
	var sb strings.Builder
	sb.WriteString("Row: ")
	for _, col := range row {
		sb.WriteString(fmt.Sprintf("%s, ", col))
	}
	p.logger.Debug(sb.String())

	sb.Reset()
	sb.WriteString("Headers: ")
	for _, header := range p.phoneFields.WriteHeaders {
		sb.WriteString(fmt.Sprintf("%s, ", header))
	}
	p.logger.Debug(sb.String())
	p.logger.Debugf("Row length: %d", len(row))
	p.logger.Debugf("Header length: %d", len(p.phoneFields.WriteHeaders))

	if len(row) != len(p.phoneFields.WriteHeaders) {
		return fmt.Errorf("row length does not match header length")
	}
	return p.writer.Write(row)
}
func (p *CSVProcessor) GetNextRow() ([]string, error) {

	row, err := p.reader.Read()

	if err != nil {
		if err.Error() == "EOF" {
			return nil, nil
		}
		return nil, err
	}

	return row, err
}

func (p *CSVProcessor) Close() error {
	p.writer.Flush()
	p.inputFile.Close()
	return p.outputFile.Close()
}
