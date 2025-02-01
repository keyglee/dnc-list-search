package dnc

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
)

type Client struct {
	dncList *os.File
	logger  *logrus.Logger
}

func (c *Client) LoadNewList(filePath string) error {

	file, err := os.Open(filePath)

	if err != nil {
		return err
	}

	defer file.Close()

	c.dncList.Truncate(0)

	buf := make([]byte, 1024)
	for {
		n, err := file.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if n == 0 {
			break
		}
		_, err = c.dncList.Write(buf[:n])
		if err != nil {
			return err
		}
	}

	return nil
}

func NewClient(dncList *os.File, logger *logrus.Logger) (*Client, error) {

	if dncList == nil {
		dncList, err := os.Open("dnc.txt")

		if err != nil {
			return nil, err
		}

		return &Client{
			dncList: dncList,
			logger:  logger,
		}, nil
	}

	return &Client{
		dncList: dncList,
		logger:  logger,
	}, nil
}

func (c *Client) FormatPhoneNumber(input string) (string, error) {
	// Extract last 10 digits
	re := regexp.MustCompile("[0-9]+")
	matches := re.FindAllString(input, -1)
	match := strings.Join(matches, "")

	c.logger.Debugf("Matched phone number: %s", match)

	if match == "" || len(match) != 10 {
		errMesg := fmt.Sprintf("Invalid phone number: %s", input)
		return "", errors.New(errMesg)
	}

	// Insert comma after 3rd digit
	return match[:3] + "," + match[3:], nil
}

func (c *Client) binarySearchFile(target string) (bool, error) {
	stat, err := c.dncList.Stat()
	if err != nil {
		return false, err
	}

	if stat.Size() == 0 {
		return false, nil
	}

	start, end := int64(0), stat.Size()
	reader := bufio.NewReader(c.dncList)

	traversals := 0

	for start <= end {
		traversals++
		mid := (start + end) / 2

		_, err := c.dncList.Seek(mid, io.SeekStart)
		if err != nil {
			c.logger.Debugf("Traversals: %d", traversals)
			return false, err
		}

		// Clear partial line if not at start
		if mid > 0 {
			reader.ReadBytes('\n')
		}

		// Read full line
		line, err := reader.ReadBytes('\n')
		c.logger.Debug(string(line))
		if err != nil {
			c.logger.Debugf("Traversals: %d", traversals)
			if err.Error() == "EOF" {
				return false, nil
			}
			return false, err
		}

		// Clean and compare
		lineStr := strings.TrimSpace(string(line))
		if lineStr == target {
			c.logger.Debugf("Traversals: %d", traversals)
			return true, nil
		}

		if lineStr < target {
			start = mid + 1
		} else {
			end = mid - 1
		}

		reader.Reset(c.dncList)
	}

	c.logger.Debugf("Traversals: %d", traversals)

	return false, nil
}

func (c *Client) Search(phoneNumber string) (bool, error) {
	// Search the phone number in the dnc list

	return c.binarySearchFile(phoneNumber)
}
