package timecode

// CountDroppedFrames counts the number of frames that have been dropped after a certain number of minutes.
// This is used when we are going the opposite direction of the below function. When this function is used,
// we're parsing a timecode into a specific frame index. We need to know how many apparent frames have been
// artificially inserted to the timecode count in order to properly subtract them out.
func CountDroppedFrames(minutes int64) int64 {

	// How many of these minutes are NOT multiples of 10?
	notMultTen := minutes - (minutes / 10)

	// Drop 2 frames for each minute that was crossed that was not a multiple of 10
	return notMultTen * 2

	// Example: Counting how many minutes, for a given input, are not multiples of 10
	// minutes => not_mult_10
	// 0 => 0
	// 1 => 1
	// 2 => 2
	// 3 => 3
	// 4 => 4
	// 5 => 5
	// 6 => 6
	// 7 => 7
	// 8 => 8
	// 9 => 9
	// 10 => 9
	// 11 => 10

}

// CountFramesToDrop counts the number of frames that should be skipped / dropped when converting the
// provided frames count into a timecode string. The result of this function is added to the actual frame count
// to produce a simulated frame count, which is them converted to a timecode string
func CountFramesToDrop(frames int64, rateNum int64) int64 {

	// Track the number of minutes that have been crossed
	var minutes int64

	// Track the number of dropped frames
	var droppedFrames int64

	// Loop until we run out of frames to count
	for frames > 0 {

		// Subtract the frames in a single minute, and bail out here if we run out
		frames -= rateNum * 60
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
