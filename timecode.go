package timecode

import (
	"fmt"
	"math"
	"time"
)

type Components struct {
	Hours, Minutes, Seconds, Frames int64
}

func (c Components) Equals(other Components) bool {
	return c.Hours == other.Hours &&
		c.Minutes == other.Minutes &&
		c.Seconds == other.Seconds &&
		c.Frames == other.Frames
}

// Timecode represents a timecode value, either as a duration or a specific point in time
type Timecode struct {
	frame int64
	rate  Rate
}

// Frame gets the frame index for this timecode
func (t *Timecode) Frame() int64 {
	return t.frame
}

// Components gets the components of the timecode: hours, minutes, seconds, frames.
func (t *Timecode) Components() Components {

	// Track the total number of frames in the timecode. If it's a drop frame rate, we need to
	// increment the number of frames to make the rest of the calculations work out.
	totalFrames := t.frame
	if t.rate.DropFrame {
		totalFrames += CountFramesToDrop(totalFrames, t.rate.Num)
	}

	// Track the remaining frames
	frames := totalFrames % int64(t.rate.Num)

	// Count the number of seconds
	totalSeconds := (totalFrames - frames) / int64(t.rate.Num)
	seconds := totalSeconds % 60

	// Count the number of minutes
	totalMinutes := (totalSeconds - seconds) / 60
	minutes := totalMinutes % 60

	// Count the total hours
	hours := (totalMinutes - minutes) / 60

	// Return the components
	return Components{
		hours,
		minutes,
		seconds,
		frames,
	}

}

// String creates a string representation for the timecode
func (t *Timecode) String() string {

	// Get the components of the timecode
	components := t.Components()

	// Determine the separator
	sep := ":"
	if t.rate.DropFrame {
		sep = ";"
	}

	// Determine the number of digits in the frame rate
	frameDigits := int(math.Log10(float64(t.rate.Num-1))) + 1

	// Create the format string for the frames
	frameFormat := fmt.Sprintf("%%0%dd", frameDigits)

	// Format the timecode
	return fmt.Sprintf(
		"%02d:%02d:%02d%s%s",
		components.Hours,
		components.Minutes,
		components.Seconds,
		sep,
		fmt.Sprintf(frameFormat, components.Frames),
	)

}

// Equals checks if this timecode is equal to another framer
func (t *Timecode) Equals(other Framer) bool {
	return other.Frame() == t.frame
}

// Add adds another framer instance to this timecode
func (t *Timecode) Add(other Framer) *Timecode {
	return &Timecode{
		frame: t.frame + other.Frame(),
		rate:  t.rate,
	}
}

// PresentationTime gets the actual presentation time of the timecode. With drop frame, this will drift from the timecode
// time before snapping back into place periodically.
func (t *Timecode) PresentationTime() time.Duration {
	return t.rate.PlaybackFrameDuration() * time.Duration(t.frame)
}
