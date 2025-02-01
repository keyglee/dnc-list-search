package dnc_test

import (
	"bufio"
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
