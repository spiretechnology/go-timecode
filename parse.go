package timecode

import (
	"errors"
	"regexp"
	"strconv"
)

// TimecodeRegex is the pattern for a valid SMPTE timecode
var TimecodeRegex = regexp.MustCompile(`^(\d\d)(:|;)(\d\d)(:|;)(\d\d)(:|;)(\d+)$`)

// MustParse parses a timecode from a string, and treats it using the provided frame rate value
func MustParse(timecode string, rate Rate) *Timecode {
	tc, err := Parse(timecode, rate)
	if err != nil {
		panic(err)
	}
	return tc
}

// Parse parses a timecode from a string, and treats it using the provided frame rate value
func Parse(timecode string, rate Rate) (*Timecode, error) {
	// Match it against the regular expression
	match := TimecodeRegex.FindStringSubmatch(timecode)
	if match == nil {
		return nil, errors.New("invalid timecode format")
	}

	// Get the components
	hours, _ := strconv.ParseInt(match[1], 10, 64)
	minutes, _ := strconv.ParseInt(match[3], 10, 64)
	seconds, _ := strconv.ParseInt(match[5], 10, 64)
	frames, _ := strconv.ParseInt(match[7], 10, 64)

	// Determine drop frame based on the final separator
	dropFrame := match[6] == ";"

	// Combine the components
	return FromComponents(Components{
		hours,
		minutes,
		seconds,
		frames,
	}, rate, dropFrame), nil
}

func FromComponents(components Components, rate Rate, dropFrame bool) *Timecode {
	// If the rate is drop frame, we need to check that the provided frame
	// isn't a dropped frame, which needs to be rounded to the nearest
	// valid frame timecode
	if dropFrame && (components.Minutes%10 > 0) && (components.Seconds == 0) && (components.Frames < rate.Drop) {
		// Move to the next valid frame in sequence
		components.Frames = rate.Drop
	}

	// Count the total number of frames in the timecode
	totalMinutes := components.Hours*60 + components.Minutes
	totalFrames := (totalMinutes*60+components.Seconds)*rate.Nominal + components.Frames

	// If it's drop frame, account for the drop incidents
	if dropFrame {
		dropFrameIncidents := totalMinutes - totalMinutes/10
		if dropFrameIncidents > 0 {
			totalFrames -= dropFrameIncidents * rate.Drop
		}
	}

	// Return the timecode with the total frames
	return &Timecode{
		frame:     totalFrames,
		rate:      rate,
		dropFrame: dropFrame,
	}
}

func FromFrame(frame int64, rate Rate, dropFrame bool) *Timecode {
	return &Timecode{
		frame,
		rate,
		dropFrame,
	}
}
