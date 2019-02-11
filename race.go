package erg

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
	BoatType         int
	Name             string //16 char limit
	Distance         int    // in Meters
	EnableStrokeData int    // 1 = yes, 0 = no
	SplitDistance    int    // Split Distance in Meters
	SplitTimes       int    // Split Times in Seconds
	Boats            []Boat
}
