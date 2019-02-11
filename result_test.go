package erg

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

// type Result struct {
// 	Place    int
// 	Time     time.Duration
// 	AvgPace  time.Duration
// 	Distance int
// 	Name     string
// 	BibNum   int
// 	Class    string
// }

func mustParse(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		panic(err)
	}
	return d
}
func TestReadResults(t *testing.T) {
	resultString := `Race Results
103

Place,Time Rowed,Meters Rowed,Boat/Team Name,Avg. Pace,ID,Bib Number,Class
1,06:27.4,2000,Kilcoyne  P.,  1:36.9,,91,
2,06:35.1,2000,Gifford  G.,  1:38.8,,83,
3,06:36.0,2000,Stahl  D.,  1:39.0,,80,

Detailed Results

Place,Lane,Rower Name,,Time,Distance,Avg. Pace,ID,Class

1,5,Kilcoyne  P.,,06:27.4,2000,1:36.9,,
2,6,Gifford  G.,,06:35.1,2000,1:38.8,,
3,7,Stahl  D.,,06:36.0,2000,1:39.0,,`

	wrongVersion := `Race Results
102

Place,Time Rowed,Meters Rowed,Boat/Team Name,Avg. Pace,ID,Bib Number,Class
1,06:27.4,2000,Kilcoyne  P.,  1:36.9,,91,
2,06:35.1,2000,Gifford  G.,  1:38.8,,83,
3,06:36.0,2000,Stahl  D.,  1:39.0,,80,

Detailed Results

Place,Lane,Rower Name,,Time,Distance,Avg. Pace,ID,Class

1,5,Kilcoyne  P.,,06:27.4,2000,1:36.9,,
2,6,Gifford  G.,,06:35.1,2000,1:38.8,,
3,7,Stahl  D.,,06:36.0,2000,1:39.0,,`

	resultsData := []Result{
		{Place: 1,
			Time:     mustParse("6m27.4s"),
			AvgPace:  mustParse("1m36.9s"),
			Distance: 2000,
			Name:     "Kilcoyne  P.",
			BibNum:   91},
		{Place: 2,
			Time:     mustParse("6m35.1s"),
			AvgPace:  mustParse("1m38.8s"),
			Distance: 2000,
			Name:     "Gifford  G.",
			BibNum:   83},
		{Place: 3,
			Time:     mustParse("6m36.0s"),
			AvgPace:  mustParse("1m39.0s"),
			Distance: 2000,
			Name:     "Stahl  D.",
			BibNum:   80}}

	reader := strings.NewReader(resultString)
	var results []Result

	err := ReadResults(&results, reader)

	if err != nil {
		t.Error("Unexpected error:", err)
	}

	if len(results) != 3 {
		t.Error("Expected 3 results but got ", len(results))
	}

	for i := range resultsData {
		if !reflect.DeepEqual(resultsData[i], results[i]) {
			t.Errorf("Expected %v but got %v", resultsData[i], results[i])
		}
	}

	// Make sure ReadResults is appending to the passed array, not overwriting
	reader = strings.NewReader(resultString)
	err = ReadResults(&results, reader)
	if err != nil {
		t.Error("Unexpected error:", err)
	}
	if len(results) != 6 {
		t.Error("expected 6 results but got", len(results))
	}

	// we should not try to read other versions
	reader = strings.NewReader(wrongVersion)
	err = ReadResults(&results, reader)
	if strings.Index(err.Error(), "version") < 0 {
		t.Error("expected wrong version but got: ", err)
	}
}

func TestReadResultsFromFile(t *testing.T) {
	results, err := ReadResultsFromFile("result_test.txt")
	if err != nil {
		t.Error("expected to read result from file but go:", err)
	}

	if len(results) != 12 {
		t.Error("expected 12 results but got", len(results))
	}

	results, err = ReadResultsFromFile("sdfsdfsdf.txt")
	if err == nil {
		t.Error("expected error but go nil")
	}

}
