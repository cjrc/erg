package erg

import (
	"fmt"
	"io"
	"os"
)

// FILESIG Signature used by Concept2 race files
const FILESIG = "RACE"

// FILEVER is the version of race file we are generating
const FILEVER = "107"

// represents BoatType in the Race structure
const (
	SINGLES = 0
	DOUBLES = 1
	FOURS   = 2
	EIGHTS  = 3
)

// Race is a flight of boats racing together, they are written to a Concept-2 .RAC file
// and imported into the Venue Racing software
type Race struct {
	BoatType         uint   // one of the consts above
	Name             string //16 char limit
	Distance         uint   // in Meters
	EnableStrokeData uint   // 1 = yes, 0 = no
	SplitDistance    uint   // Split Distance in Meters
	SplitTime        uint   // Split Time in Seconds
	Boats            []Boat // len(boats) does not need to equal NLanes
	NLanes           uint   // Number of lanes in this race
	DurationType     uint   // 0=distance, 1=time?
}

// Given a list of boats from a race, this will return the boat that is
// in the specified lane.  Returns an empty boat if that lane is empty
func findByLane(boats []Boat, lane uint) Boat {
	// An empty lane must still have a non-empty name
	// for the venue software to accept the race file
	empty := Boat{Name: " "}

	for _, boat := range boats {
		if lane == boat.Lane {
			return boat
		}
	}

	return empty
}

// Write the race in the format understood by the Concept 2 Venue software
// TODO:
// These generated race files will not currently work
// with team boats
// see https://c2usa.fogbugz.com/?W44
func (race Race) Write(w io.Writer) error {
	// The race won't start if the name is longer than 16 characters
	maxNameLen := 16
	if len(race.Name) < maxNameLen {
		maxNameLen = len(race.Name)
	}

	// file type signature
	if _, err := fmt.Fprintln(w, FILESIG); err != nil {
		return err
	}

	// file format ver 107 including class.
	if _, err := fmt.Fprintln(w, FILEVER); err != nil {
		return err
	}
	// file type signature
	// file format ver 107 including class.
	// team config (singles=0, doubles=1, fours=2, eights=3)
	// race name; see note 1
	// distance in meters
	// Duration Type
	// Next line is always 0 (was View Mode in older days)
	if _, err := fmt.Fprintf(w, "%s\n%s\n%d\n%s\n%d\n%d\n0\n",
		FILESIG, FILEVER, race.BoatType, race.Name[:maxNameLen],
		race.Distance, race.DurationType); err != nil {
		return err
	}

	// EnableStrokeData 1=yes,0=no
	// SplitDistance in Meters
	// SplitTimes in Seconds
	// Actual Number of Boats in this race (2-40)
	if _, err := fmt.Fprintf(w, "%d\n%d\n%d\n%d\n",
		race.EnableStrokeData, race.SplitDistance,
		race.SplitTime, race.NLanes); err != nil {
		return err
	}

	for lane := uint(1); lane <= race.NLanes; lane++ {
		boat := findByLane(race.Boats, lane)
		if err := boat.Write(w); err != nil {
			return err
		}
	}

	// Concept 2 example file has this closing 0
	// It's always 0 for individual races
	// should be total number of PMs for team races
	_, err := fmt.Fprintln(w, "0")

	return err
}

// WriteToFile save race in the specified file in a format
// that is readable by the Venue Racing Software
func (race Race) WriteToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return race.Write(file)
}
