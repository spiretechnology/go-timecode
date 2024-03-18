package timecode

import (
	"math"
	"strconv"
	"strings"
)

const (
	dropOccurrencesPerHour = 54
	secondsPerHour         = 60 * 60
)

var (
	Rate_23_976 = Rate{"23.976", 24, 0, 24000, 1001}
	Rate_24     = Rate{"24", 24, 0, 24, 1}
	Rate_30     = Rate{"30", 30, 0, 30, 1}
	Rate_29_97  = Rate{"29.97", 30, 2, 30000, 1001}
	Rate_60     = Rate{"60", 60, 0, 60, 1}
	Rate_59_94  = Rate{"59.94", 60, 4, 60000, 1001}
)

// Rate represents a frame rate for a timecode
type Rate struct {
	Str      string
	Nominal  int
	Drop     int
	Num, Den int
}

// String creates a string representation of the rate
func (r *Rate) String() string {
	return r.Str
}

// ParseRate returns a Rate from a string representation.
func ParseRate(str string) (Rate, bool) {
	switch str {
	case "23.976", "23.98":
		return Rate_23_976, true
	case "24":
		return Rate_24, true
	case "30":
		return Rate_30, true
	case "29.97":
		return Rate_29_97, true
	case "60":
		return Rate_60, true
	case "59.94":
		return Rate_59_94, true
	}
	return Rate{}, false
}

// RateFromFraction returns a Rate from a numerator and denominator.
func RateFromFraction(num, den int) Rate {
	type fraction struct {
		num, den int
	}
	switch (fraction{num, den}) {
	case fraction{24000, 1001}:
		return Rate_23_976
	case fraction{24, 1}:
		return Rate_24
	case fraction{30, 1}:
		return Rate_30
	case fraction{30000, 1001}:
		return Rate_29_97
	case fraction{60, 1}:
		return Rate_60
	case fraction{60000, 1001}:
		return Rate_59_94
	}

	// Calculate the nominal frame rate (number of frames in a second without drops)
	var nominal int
	if num%den == 0 {
		nominal = num / den
	} else {
		nominal = den - (num % den)
	}

	// Format it as a string (ie. 23.976)
	str := strconv.FormatFloat(float64(num)/float64(den), 'f', 3, 64)
	for strings.HasSuffix(str, "0") {
		str = str[:len(str)-1]
	}
	str = strings.TrimSuffix(str, ".")

	// Calculate the number of frames to skip per drop occurrence
	actualFramesPerHour := num * secondsPerHour / den
	nominalFramesPerHour := nominal * secondsPerHour
	totalFramesDropped := nominalFramesPerHour - actualFramesPerHour
	framesPerDrop := int(math.Round(float64(totalFramesDropped) / float64(dropOccurrencesPerHour)))

	return Rate{
		Str:     str,
		Nominal: nominal,
		Drop:    framesPerDrop,
		Num:     num,
		Den:     den,
	}
}
