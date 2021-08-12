package timecode

import (
	"errors"
	"regexp"
	"strconv"
	"time"
)

// TimecodeRegex is the pattern for a valid SMPTE timecode
var TimecodeRegex = regexp.MustCompile(`^(\d\d)(:|;)(\d\d)(:|;)(\d\d)(:|;)(\d\d)$`)

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

	// Combine the components
	return FromComponents(Components{
		hours,
		minutes,
		seconds,
		frames,
	}, rate)

}

func FromComponents(components Components, rate Rate) (*Timecode, error) {

	// If the rate is drop frame, we need to check that the provided frame
	// isn't a dropped frame, which needs to be rounded to the nearest
	// valid frame timecode
	if rate.DropFrame {

		// If it's a dropped frame
		if (components.Minutes%10 > 0) && (components.Seconds == 0) && (components.Frames == 0 || components.Frames == 1) {

			// Move to the next valid frame in sequence
			components.Frames = 2

		}

	}

	// Count up the total number of frames
	totalSeconds := (((components.Hours * 60) + components.Minutes) * 60) + components.Seconds
	totalFrames := totalSeconds*rate.Num + components.Frames

	// Count the number of dropped frames
	if rate.DropFrame {
		totalFrames -= CountDroppedFrames(components.Minutes)
	}

	// Return the timecode with the total frames
	return &Timecode{
		frame: totalFrames,
		rate:  rate,
	}, nil

}

func FromFrame(frame int64, rate Rate) *Timecode {
	return &Timecode{
		frame,
		rate,
	}
}

func FromPresentationTime(presentationTime time.Duration, rate Rate) *Timecode {
	return &Timecode{
		frame: int64(presentationTime / rate.PlaybackFrameDuration()),
		rate:  rate,
	}
}
