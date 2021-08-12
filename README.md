# go-timecode

A Go library for parsing and manipulating SMPTE timecodes and frame rates.

Development of this library is very test-driven to ensure accuracy of the frame and timecode calculations. Adding additional test cases is encouraged!

## Parse a timecode

### Example with drop frame (23.976 FPS)
```go
tc, err := timecode.Parse("00:01:02;23", timecode.Rate_23_976)
tc.String() // => 00:01:02;23
tc.Frame() // => 1509
```

### Example without drop frame (24 FPS)
```go
tc, err := timecode.Parse("00:01:02:23", timecode.Rate_24)
tc.String() // => 00:01:02:23
tc.Frame() // => 1511
```

