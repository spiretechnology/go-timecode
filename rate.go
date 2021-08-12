package timecode

import (
	"fmt"
)

var (
	// 23.976 has exactly 24 frames per second. However, the textual representation of timecodes using this rate
	// skip two frames every minute, except when the minute is a multiple of 10. This is because 23.976 footage
	// actually does display at a rate of 23.976 frames each second on televisions. To ensure that the first
	// timecode in an hour of footage is 00:00:00;00 and the last timecode in that hour is 01:00:00;00, drop
	// frame was invented. It is purely a matter of presentation.
	Rate_23_976 = Rate{24, true}

	// Standard 24 FPS, with no drop frame
	Rate_24 = Rate{24, false}

	// Other formats...
	Rate_30    = Rate{30, false}
	Rate_29_97 = Rate{30, true}
	Rate_60    = Rate{60, false}
	Rate_59_94 = Rate{60, true}
)

// Rate represents a frame rate for a timecode
type Rate struct {
	Num       int64
	DropFrame bool
}

// String creates a string representation of the rate
func (r *Rate) String() string {
	return fmt.Sprintf("%d", r.Num)
}
