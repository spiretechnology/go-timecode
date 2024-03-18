//go:build exhaustivetest

package timecode_test

import (
	"bufio"
	"os"
	"path/filepath"
	"testing"

	"github.com/spiretechnology/go-timecode"
	"github.com/stretchr/testify/require"
)

func runTimecodesTest(t *testing.T, rate timecode.Rate, testfile string) {
	file, err := os.Open(filepath.Join("testdata", testfile))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var prevTimecode *timecode.Timecode

	for frameIndex := int64(0); scanner.Scan(); frameIndex++ {
		line := scanner.Text()
		if line == "" {
			continue
		}

		// Frame index -> timecode string
		tcFromIndex := timecode.FromFrame(frameIndex, rate, rate.Drop != 0)
		require.Equal(t, line, tcFromIndex.String(), "frame %d", frameIndex)

		// Timecode string -> frame index
		tcFromStr, err := timecode.Parse(line, rate)
		require.NoError(t, err, "error parsing timecode %s", line)
		require.Equal(t, frameIndex, tcFromStr.Frame(), "timecode %s", line)

		// Compare to previous timecode
		if prevTimecode != nil {
			prevPlusOne := prevTimecode.AddFrames(1)
			require.Equal(t, tcFromStr.String(), prevPlusOne.String(), "timecode %s", line)
			require.Equal(t, frameIndex, prevPlusOne.Frame(), "timecode %s", line)
		}

		prevTimecode = tcFromStr
	}
}

func TestAllTimecodes(t *testing.T) {
	t.Run("all timecodes - 23.976", func(t *testing.T) {
		runTimecodesTest(t, timecode.Rate_23_976, "tc-all-23_976.txt")
	})
	t.Run("all timecodes - 24", func(t *testing.T) {
		runTimecodesTest(t, timecode.Rate_24, "tc-all-24.txt")
	})
	t.Run("all timecodes - 29.97", func(t *testing.T) {
		runTimecodesTest(t, timecode.Rate_29_97, "tc-all-29_97.txt")
	})
	t.Run("all timecodes - 30", func(t *testing.T) {
		runTimecodesTest(t, timecode.Rate_30, "tc-all-30.txt")
	})
	t.Run("all timecodes - 59.94", func(t *testing.T) {
		runTimecodesTest(t, timecode.Rate_59_94, "tc-all-59_94.txt")
	})
	t.Run("all timecodes - 60", func(t *testing.T) {
		runTimecodesTest(t, timecode.Rate_60, "tc-all-60.txt")
	})
}
