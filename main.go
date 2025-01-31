package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
)

func LoadNewList(filePath string) ([]string, error) {

	return nil, nil
}

func binarySearchFile(file *os.File, target string, logger *logrus.Logger) (bool, error) {
	stat, err := file.Stat()
	if err != nil {
		return false, err
	}

	start, end := int64(0), stat.Size()
	reader := bufio.NewReader(file)

	traversals := 0

	for start <= end {
		traversals++
		mid := (start + end) / 2

		_, err := file.Seek(mid, 0)
		if err != nil {
			logger.Debugf("Traversals: %d", traversals)
			return false, err
		}

		// Clear partial line if not at start
		if mid > 0 {
			reader.ReadBytes('\n')
		}

		// Read full line
		line, err := reader.ReadBytes('\n')
		logger.Debug(string(line))
		if err != nil {
			logger.Debugf("Traversals: %d", traversals)
			if err.Error() == "EOF" {
				return false, nil
			}
			return false, err
		}

		// Clean and compare
		lineStr := strings.TrimSpace(string(line))
		if lineStr == target {
			logger.Debugf("Traversals: %d", traversals)
			return true, nil
		}

		if lineStr < target {
			start = mid + 1
		} else {
			end = mid - 1
		}

		reader.Reset(file)
	}

	logger.Debugf("Traversals: %d", traversals)

	return false, nil
}

func formatPhoneNumber(input string, logger *logrus.Logger) (string, error) {
	// Extract last 10 digits
	re := regexp.MustCompile("[0-9]+")
	matches := re.FindAllString(input, -1)
	match := strings.Join(matches, "")

	logger.Debugf("Matched phone number: %s", match)

	if match == "" || len(match) != 10 {
		errMesg := fmt.Sprintf("Invalid phone number: %s", input)
		return "", errors.New(errMesg)
	}

	// Insert comma after 3rd digit
	return match[:3] + "," + match[3:], nil
}

func initLogger(logLevel string) *logrus.Logger {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		level = logrus.ErrorLevel // Set default if invalid
	}
	logger.SetLevel(level)

	// Optional: Add timestamp and caller info
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	return logger
}

func main() {
	loadNewCard := flag.String("new", "", "Adds a new card from the text file")
	logLevel := flag.String("log", "error", "Log level, can be error, info, or debug")
	flag.Parse()

	logger := initLogger(*logLevel)

	file, err := os.Open("dnc.txt")

	if err != nil {
		logger.Error(loadNewCard)
		panic(err)
	}

	args := flag.Args()

	for _, arg := range args {

		logger.Debugf("Processing argument: '%s'", arg)
		// Format phone number
		formattedArg, err := formatPhoneNumber(arg, logger)

		logger.Debugf("Formatted argument: '%s'", formattedArg)
		if err != nil {
			logger.Error(err)
			panic(err)
		}
		// Argument needs to be in the form of 'xxx,xxxxxxx'
		found, err := binarySearchFile(file, formattedArg, logger)

		if err != nil {
			panic(err)
		}

		if found {
			fmt.Printf("%s Found in the file\n", arg)
		} else {
			fmt.Printf("%s Not Found in the file\n", arg)
		}
	}
}
