package erg

import (
	"fmt"
	"io"
)

// Boat represents a single entry in a race
type Boat struct {
	Name    string // participant name; see note 2 (avoid punctuation)
	ID      uint   // Bib Number
	Class   string // Class; see note 3 (fine to leave this blank)
	Country string // Country Code; see note 4 (set this to "USA")
	DOB     string // Can be left blank, format is MMDDYYYY
}

// Write outputs a boat in the right format for a .RAC file
func (boat Boat) Write(w io.Writer) {
	name := boat.Name
	if name == "" {
		name = " "
	}
	fmt.Fprintln(w, name)
	fmt.Fprintln(w, boat.ID)
	fmt.Fprintln(w, boat.Class)
	fmt.Fprintln(w, boat.Country)
	fmt.Fprintln(w, boat.DOB)
}
