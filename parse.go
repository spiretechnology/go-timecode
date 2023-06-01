package timecode

import (
	"errors"
	"regexp"
	"strconv"
)

// TimecodeRegex is the pattern for a valid SMPTE timecode
var TimecodeRegex = regexp.MustCompile(`^(\d\d)(:|;)(\d\d)(:|;)(\d\d)(:|;)(\d+)$`)

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
	}, rate, dropFrame)
}

func FromComponents(components Components, rate Rate, dropFrame bool) (*Timecode, error) {
	// If the rate is drop frame, we need to check that the provided frame
	// isn't a dropped frame, which needs to be rounded to the nearest
	// valid frame timecode
	if dropFrame && (components.Minutes%10 > 0) && (components.Seconds == 0) && (components.Frames == 0 || components.Frames == 1) {
		// Move to the next valid frame in sequence
		components.Frames = 2
	}

	// Count the total number of frames in the timecode, ignoring dropped frames for now
	totalFrames := ((((components.Hours*60)+components.Minutes)*60)+components.Seconds)*rate.RoundUp() + components.Frames

	// If we're in drop frame, count the number of frames that need to be subtracted from the frame count.
	// This is equal to the number of 1-minute intervals that are not multiples of 10, times 2.
	if dropFrame {
		totalFrames -= (components.Minutes - components.Minutes/10) * 2
	}

	// Return the timecode with the total frames
	return &Timecode{
		frame:     totalFrames,
		rate:      rate,
		dropFrame: dropFrame,
	}, nil
}

func FromFrame(frame int64, rate Rate, dropFrame bool) *Timecode {
	return &Timecode{
		frame,
		rate,
		dropFrame,
	}
}
