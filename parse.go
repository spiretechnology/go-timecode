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
	return FromComponents(&Components{
		hours,
		minutes,
		seconds,
		frames,
	}, rate)

}

func FromComponents(components *Components, rate Rate) (*Timecode, error) {

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

func FromFrames(frames int64, rate Rate) *Timecode {
	return &Timecode{
		frame: frames,
		rate:  rate,
	}
}

func FromPresentationTime(presentationTime time.Duration, rate Rate) *Timecode {
	return &Timecode{
		frame: int64(presentationTime / rate.PlaybackFrameDuration()),
		rate:  rate,
	}
}
