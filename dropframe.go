package timecode

// CountFramesToDrop counts the number of frames that should be skipped / dropped when converting the
// provided frames count into a timecode string. The result of this function is added to the actual frame count
// to produce a simulated frame count, which is them converted to a timecode string
func CountFramesToDrop(frames int64, rate Rate) int64 {
	// Track the number of minutes that have been crossed, and the number of frames dropped so far
	var minutes, droppedFrames int64

	// Loop until we run out of frames to count
	rateRoundedUp := rate.RoundUp()
	for frames > 0 {
		// Subtract the frames in a single minute, and bail out here if we run out
		frames -= rateRoundedUp * 60
		if frames < 0 {
			break
		}

		// Increment the number of minutes
		minutes++

		// If this minute is not a multiple of 10, it will include two dropped frames
		if minutes%10 > 0 {
			// Track these dropped frames
			droppedFrames += 2

			// Add the frames to the frames count, extending our work by 2 additional frames. We do this to ensure that
			// every minute-crossing is tracked, as dropped frames can create additional minute-crossings that would never
			// have taken place without the artificial inflation of the dropped frames.
			frames += 2
		}
	}

	// Return the total dropped frames
	return droppedFrames
}
