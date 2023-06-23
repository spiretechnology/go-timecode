package timecode

import (
	"fmt"
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
	frame     int64
	rate      Rate
	dropFrame bool
}

// Frame gets the frame index for this timecode
func (t *Timecode) Frame() int64 {
	return t.frame
}

func (t *Timecode) componentsNDF(frame int64) Components {
	// Track the remaining frames
	frames := frame % t.rate.Nominal

	// Count the number of seconds
	totalSeconds := (frame - frames) / t.rate.Nominal
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

func (t *Timecode) componentsDF(frame int64) Components {
	// Calculate the NDF components
	comps := t.componentsNDF(frame)

	// Count the total number of minutes crossed
	minutesCrossed := comps.Hours*60 + comps.Minutes
	dropFrameIncidents := minutesCrossed - minutesCrossed/10

	// As long as there are unhandled drop frame incidents
	for dropFrameIncidents > 0 {
		// Add the appropriate number of frames
		frame += dropFrameIncidents * t.rate.Drop

		// Recalculate the NDF components
		newComps := t.componentsNDF(frame)

		// Count the number of drop frame incidents
		dropFrameIncidents = 0
		for m := comps.Hours*60 + comps.Minutes + 1; m <= newComps.Hours*60+newComps.Minutes; m++ {
			if m%10 > 0 {
				dropFrameIncidents++
			}
		}

		// Update the components
		comps = newComps
	}
	return comps
}

// Components gets the components of the timecode: hours, minutes, seconds, frames.
func (t *Timecode) Components() Components {
	if !t.dropFrame {
		return t.componentsNDF(t.frame)
	} else {
		return t.componentsDF(t.frame)
	}
}

// String creates a string representation for the timecode
func (t *Timecode) String() string {
	// Get the components of the timecode
	components := t.Components()

	// Determine the separator
	sep := ":"
	if t.dropFrame {
		sep = ";"
	}

	// Determine the number of digits in the frame rate, and create the format string. We do this to account
	// for triple-digit frame rates.
	frameDigits := len(fmt.Sprintf("%d", t.rate.Nominal))
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
		frame:     t.frame + other.Frame(),
		rate:      t.rate,
		dropFrame: t.dropFrame,
	}
}

func (t *Timecode) AddFrames(other int64) *Timecode {
	return t.Add(Frame(other))
}
