package erg

import (
	"fmt"
	"io"
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
func (race Race) Write(w io.Writer) {
	// The race won't start if the name is longer than 16 characters
	maxNameLen := 16
	if len(race.Name) < maxNameLen {
		maxNameLen = len(race.Name)
	}

	fmt.Fprintln(w, FILESIG)                // file type signature
	fmt.Fprintln(w, FILEVER)                // file format ver 107 including class.
	fmt.Fprintln(w, race.BoatType)          // team config (singles=0, doubles=1, fours=2, eights=3)
	fmt.Fprintln(w, race.Name[:maxNameLen]) // race name; see note 1
	fmt.Fprintln(w, race.Distance)          // distance in meters
	fmt.Fprintln(w, "0")                    // duration type is distance
	fmt.Fprintln(w, "0")                    // this line is always zero
	fmt.Fprintln(w, race.EnableStrokeData)  // 1 = yes, 0 = no
	fmt.Fprintln(w, race.SplitDistance)     // SplitDistance in Meters
	fmt.Fprintln(w, race.SplitTime)         // Split Times in Seconds
	fmt.Fprintln(w, race.NLanes)            // Actual Number of boats in this race (2-40)

	for lane := uint(1); lane <= race.NLanes; lane++ {
		boat := findByLane(race.Boats, lane)
		boat.Write(w)
	}

	fmt.Fprintln(w, "0") // example file has this closing 0
}
