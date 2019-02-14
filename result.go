package erg

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

// Result represents the race result of one erg from the Venue Racing results file
type Result struct {
	Place    int           `db:"place"`
	Time     time.Duration `db:"time"`
	AvgPace  time.Duration `db:"avg_pace"`
	Distance int           `db:"distance"`
	Name     string        `db:"name"`
	BibNum   int           `db:"bib_num"`
	Class    string        `db:"class"`
}

//ReadResults reads the race results from the specified io.Reader and appends them to the
//supplied results array.
//It will return an error if the read results are in an invalid format or are not
//Version 103 results.
func ReadResults(results *[]Result, reader io.Reader) error {
	scanner := bufio.NewScanner(reader)

	scanner.Scan()
	if scanner.Text() != "Race Results" {
		return fmt.Errorf("invalid or corrupted race results")
	}

	scanner.Scan()
	if ver := scanner.Text(); ver != "103" {
		return fmt.Errorf("found version %v results -- this software only knows version 103", ver)
	}

	scanner.Scan() // Skip blank line
	scanner.Scan() // Skip headers line

	lineNumber := 5
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break // a blank line is the end of results
		}

		parts := strings.Split(line, ",")
		if len(parts) != 8 {
			return fmt.Errorf("invalid results on line %d", lineNumber)
		}

		tmp := strings.Replace(parts[1], ":", "m", -1) + "s"
		finishTime, err := time.ParseDuration(tmp)
		if err != nil {
			return fmt.Errorf("invalid race time on line %d: %v", lineNumber, err)
		}

		tmp = strings.TrimSpace(strings.Replace(parts[4], ":", "m", -1) + "s")
		avgPace, err := time.ParseDuration(tmp)
		if err != nil {
			return fmt.Errorf("invalid average pace on line %d: %v", lineNumber, err)
		}

		distance, err := strconv.Atoi(parts[2])
		if err != nil {
			return fmt.Errorf("invalid race distance on line %d: %v", lineNumber, err)
		}

		place, err := strconv.Atoi(parts[0])
		if err != nil {
			return fmt.Errorf("invalid finish place on line %d: %v", lineNumber, err)
		}

		bibNum, err := strconv.Atoi(parts[6])
		if err != nil {
			return fmt.Errorf("invalid id on line %d: %v", lineNumber, err)
		}

		result := Result{
			Place:    place,
			Time:     finishTime,
			Distance: distance,
			Name:     parts[3],
			AvgPace:  avgPace,
			BibNum:   bibNum,
			Class:    parts[7],
		}
		*results = append(*results, result)

		lineNumber++
	}
	return scanner.Err()
}

// ReadResultsFromFile is a convenience function reads result from the specified filed
func ReadResultsFromFile(filename string) ([]Result, error) {
	results := make([]Result, 0)

	file, err := os.Open(filename)
	if err != nil {
		return results, err
	}

	defer file.Close()

	err = ReadResults(&results, file)

	return results, err
}
