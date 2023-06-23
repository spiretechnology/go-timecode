package timecode

var (
	// 23.976 has exactly 24 frames per second. However, the textual representation of timecodes using this rate
	// skip two frames every minute, except when the minute is a multiple of 10. This is because 23.976 footage
	// actually does display at a rate of 23.976 frames each second on televisions. To ensure that the first
	// timecode in an hour of footage is 00:00:00;00 and the last timecode in that hour is 01:00:00;00, drop
	// frame was invented. It is purely a matter of presentation.
	Rate_23_976 = Rate{"23.976", 24, 2, 24000, 1001}
	Rate_24     = Rate{"24", 24, 0, 24, 1}
	Rate_30     = Rate{"30", 30, 0, 30, 1}
	Rate_29_97  = Rate{"29.97", 30, 2, 30000, 1001}
	Rate_60     = Rate{"60", 60, 0, 60, 1}
	Rate_59_94  = Rate{"59.94", 60, 4, 60000, 1001}
)

// Rate represents a frame rate for a timecode
type Rate struct {
	Str      string
	Nominal  int64
	Drop     int64
	Num, Den int64
}

// String creates a string representation of the rate
func (r *Rate) String() string {
	return r.Str
}
