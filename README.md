# go-timecode

A Go library for parsing and manipulating SMPTE timecodes and frame rates.

Development of this library is very test-driven to ensure accuracy of the frame and timecode calculations. If you'd like to contribute to this library, adding additional useful test cases is a great place to start!

## Parse a timecode

#### Drop frame (23.976 FPS)
```go
tc, err := timecode.Parse("00:01:02;23", timecode.Rate_23_976)
tc.String() // => 00:01:02;23
tc.Frame() // => 1509
```

#### Non-drop frame (24 FPS)
```go
tc, err := timecode.Parse("00:01:02:23", timecode.Rate_24)
tc.String() // => 00:01:02:23
tc.Frame() // => 1511
```

## Create a timecode from a frame count
```go
tc := timecode.FromFrame(1511, timecode.Rate_24)
tc.String() // => 00:01:02:23
tc.Frame() // => 1511
```

## Algebra with timecodes and frames
```go
tc, err := timecode.Parse("00:01:02:23", timecode.Rate_24)
tc = tc.Add(timecode.Frame(3))
tc.String() // => 00:01:03:02
tc.Frame() // => 1514
```
