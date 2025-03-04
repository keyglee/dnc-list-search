package main

import (
	"dnclistsearch/csv"
	"dnclistsearch/dnc"
	"dnclistsearch/logger"
	"dnclistsearch/output"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	speedrun := flag.Bool("speedrun", false, "Runs the program with timer")
	logLevel := flag.String("log", "error", "Log level, can be error, info, or debug")
	prettyPrint := flag.Bool("pretty", false, "Pretty print the output")
	responseDelimiter := flag.String("delimiter", "none", "Delimiter for the response")
	responseSeparator := flag.String("separator", "newline", "Separator for the response")
	inputCSV := flag.String("csv", "", "Path to input CSV file")
	outputCSV := flag.String("output", "results.csv", "Path to output CSV file")

	flag.Parse()

	var startTime time.Time
	if *speedrun {
		startTime = time.Now()
	}

	logger := logger.NewLogger(*logLevel)

	dncListFile, err := os.Open("dnc.txt")

	if err != nil {
		logger.Errorf("Error opening dnc list file: %s\nIts likely you're missing the dnc.txt file in your directory\nMake sure to run this in the same directory as dnc.txt", err)
		return
	}

	dncClient, err := dnc.NewClient(dncListFile, logger)

	if err != nil {
		panic(err)
	}

	wantedHeaders := []string{"Primary #", "Phone 1", "Phone 2", "Phone 3"}

	if *inputCSV != "" {
		processor, err := csv.NewCSVProcessor(*inputCSV, *outputCSV, wantedHeaders)
		if err != nil {
			logger.Error(err)
			panic(err)
		}

		defer processor.Close()

		for {
			csvRow, err := processor.GetNextRow()

			if err != nil {
				panic(err)
			}

			if csvRow == nil {
				break
			}

			phones := processor.GetPhoneNumbers(csvRow)

			foundList := make([]string, 0)

			for _, phone := range phones {

				formatPhone, err := dncClient.FormatPhoneNumber(phone)
				if err != nil {
					foundList = append(foundList, "either no phone or invalid phone")
				} else {
					found, err := dncClient.Search(formatPhone)

					if err != nil {
						foundList = append(foundList, "encountered an error")
					} else if found {
						foundList = append(foundList, "found in DNC list")
					} else {
						foundList = append(foundList, "not found in DNC list")
					}
				}
			}

			logger.Debugf("Found list length: %v", len(foundList))
			logger.Debugf("CSV Row: %v", len(csvRow))

			csvRow = append(csvRow, foundList...)

			logger.Debugf("Writing row: %v", csvRow)

			err = processor.WriteRow(csvRow)

			if err != nil {
				panic(err)
			}
		}

	} else {

		args := flag.Args()

		var response []string

		for _, arg := range args {

			logger.Debugf("Processing argument: '%s'", arg)
			// Format phone number
			formattedArg, err := dncClient.FormatPhoneNumber(arg)

			logger.Debugf("Formatted argument: '%s'", formattedArg)
			if err != nil {
				logger.Error(err)
				panic(err)
			}
			// Argument needs to be in the form of 'xxx,xxxxxxx'
			found, err := dncClient.Search(formattedArg)

			if err != nil {
				panic(err)
			}

			deliminator := output.GetDelimiter(*responseDelimiter)

			if !*prettyPrint {
				response = append(response, fmt.Sprintf("%s%t%s", deliminator, found, deliminator))
			} else if found {
				response = append(response, fmt.Sprintf("%s%s Found in the file%s", deliminator, arg, deliminator))
			} else {
				response = append(response, fmt.Sprintf("%s%s Not Found in the file%s", deliminator, arg, deliminator))
			}
		}

		fmt.Println(strings.Join(response, string(output.GetSeparator(*responseSeparator))))
	}

	if *speedrun {
		elapsed := time.Since(startTime)
		fmt.Println("Execution time: ", elapsed)
	}
}
