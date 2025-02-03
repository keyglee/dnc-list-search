package dnc_test

import (
	"bufio"
	"dnclistsearch/csv"
	"dnclistsearch/dnc"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestDNCSearch(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	// Open DNC list file
	dncFile, err := os.Open("dnc_test.txt")
	if err != nil {

		t.Fatal(err)
	}

	// Create client
	client, err := dnc.NewClient(dncFile, logger)
	if err != nil {
		t.Fatal(err)
	}

	// Read and test first half of file

	dncTestFile, err := os.Open("dnc_test.txt")
	if err != nil {
		t.Fatal(err)
	}

	lineCount := 0

	testReader := bufio.NewReader(dncTestFile)
	// Read until the end of the file
	for {
		line, err := testReader.ReadBytes('\n')
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			t.Fatal(err)
			break
		}
		// Test line
		found, err := client.Search(string(line))
		if err != nil {
			t.Fatal(err)
			break
		}

		if found {
			t.Errorf("Failed to find %s", string(line))
			break
		} else {
			lineCount++
		}
	}

	t.Logf("Successfully found %d phone numbers", lineCount)
}

func TestCSVProcessor(t *testing.T) {
	// Mock CSV content
	csvData := `TR Sheet,Intake Rep,LLC Name,Full Name,Address,Primary #,Type,Phone 1,Type 1,Phone 2,Type 2,Phone 3,Type 3,Email
Test1,Rep1,LLC1,John Doe,123 St,(555) 123-4567,Cell,5551234568,Home,5551234569,Work,5551234570,Mobile,test@test.com
Test2,Rep2,LLC2,Jane Doe,456 St,,Cell,5551234571,Home,,,,,test2@test.com`

	// Create temp files
	tmpInput := t.TempDir() + "/test_input.csv"
	tmpOutput := t.TempDir() + "/test_output.csv"

	// Write test data
	if err := os.WriteFile(tmpInput, []byte(csvData), 0644); err != nil {
		t.Fatal(err)
	}

	// Initialize processor
	processor, err := csv.NewCSVProcessor(tmpInput, tmpOutput)
	if err != nil {
		t.Fatal(err)
	}
	defer processor.Close()

	// Read first data row
	row, err := processor.GetReader().Read()
	if err != nil {
		t.Fatal(err)
	}

	// Test phone number extraction
	phones := processor.GetPhoneNumbers(row)
	expected := []string{"(555) 123-4567", "5551234568", "5551234569", "5551234570"}

	if len(phones) != len(expected) {
		t.Errorf("Expected %d phone numbers, got %d", len(expected), len(phones))
	}

	for i, phone := range phones {
		if phone != expected[i] {
			t.Errorf("Expected phone %s, got %s", expected[i], phone)
		}
	}

	// Test row with missing phones
	row, err = processor.reader.Read()
	if err != nil {
		t.Fatal(err)
	}

	phones = processor.GetPhoneNumbers(row)
	if len(phones) != 1 || phones[0] != "5551234571" {
		t.Errorf("Expected single phone number 5551234571, got %v", phones)
	}
}
