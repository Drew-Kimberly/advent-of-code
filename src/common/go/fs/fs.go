package fs

import (
	"bufio"
	"os"
)

func ExtractInputLines(filePath string) ([]string, error) {
	inputFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(inputFile)
	scanner.Split(bufio.ScanLines)

	var inputLines []string
	for scanner.Scan() {
		inputLines = append(inputLines, scanner.Text())
	}

	err = inputFile.Close()
	if err != nil {
		return nil, err
	}

	return inputLines, nil
}
