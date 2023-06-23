# go-timecode

A Go library for parsing and manipulating SMPTE timecodes and frame rates.

Development of this library is very test-driven to ensure accuracy of the frame and timecode calculations. If you'd like to contribute to this library, adding additional useful test cases is a great place to start!

## Installation
```sh
go get github.com/spiretechnology/go-timecode
```

## Usage Examples
### Parse a timecode (drop frame)
```go
tc, err := timecode.Parse("00:01:02;23", timecode.Rate_23_976)
tc.String() // => 00:01:02;23
tc.Frame() // => 1509
```

### Parse a timecode (non-drop frame)
```go
tc, err := timecode.Parse("00:01:02:23", timecode.Rate_24)
tc.String() // => 00:01:02:23
tc.Frame() // => 1511
```

### Create a timecode from a frame count
```go
tc := timecode.FromFrame(1511, timecode.Rate_24, false /* non-drop frame */)
tc.String() // => 00:01:02:23
tc.Frame() // => 1511
```

### Algebra with timecodes and frames
```go
tc, err := timecode.Parse("00:01:02:23", timecode.Rate_24)
tc = tc.Add(timecode.Frame(3))
tc.String() // => 00:01:03:02
tc.Frame() // => 1514
```

## Note: parsing timecodes that don't exist in drop frame

Drop frame timecodes skip the first 2 frames of each minute, unless the minute is a multiple of 10. This changes to the first 4 frames of each minute if the frame rate is 59.94.

For instance, in `23.976`, the timecode `00:00:59:23` is immediately followed by `00:01:00:02`. Two timecodes were dropped: `00:01:00:00` and `00:01:00:01`

Those dropped timecodes don't correspond to any actual frame number, and so we need to choose how to resolve those frames. The choice we have made with this library is to round up the next valid frame. If you try to parse `00:01:00:00`, the result will be rounded up to `00:01:00:02`, which is the next valid frame in the sequence.

## Contributing
We welcome contributions that make this library more reliable. To add test cases, fix bugs, or anything else, please submit a pull request.

## Other resources

- [spiretechnology/js-timecode](https://github.com/spiretechnology.com/go-timecode) - A TypeScript / JavaScript library for parsing and manipulating SMPTE timecodes and frame rates.
