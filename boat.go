package erg

import (
	"fmt"
	"io"
)

// Boat represents a (single,double,four,etc..) in a Venue Race
type Boat struct {
	Name    string // participant or boat name; see note 2 (avoid punctuation)
	BibNum  uint   // Bib Number
	Class   string // Class; see note 3 (fine to leave this blank)
	Country string // Country Code; see note 4 (set this to "USA")
	DOB     string // Can be left blank, format is MMDDYYYY
	Lane    uint   // Erg Lane assignment
}

// Write outputs a boat in the right format for a .RAC file
// The lane is not written out here.. a Race will access the Lane number to
// output the boats in the right order
func (boat Boat) Write(w io.Writer) error {
	name := boat.Name
	if name == "" {
		name = " "
	}
	_, err := fmt.Fprintf(w, "%s\n%d\n%s\n%s\n%s\n",
		name, boat.BibNum, boat.Class, boat.Country, boat.DOB)
	return err
}
