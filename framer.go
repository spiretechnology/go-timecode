package timecode

// Framer defines the interface for something that can be converted into a specific
// frame index. Currently, this is limited only to Frame and Timecode
type Framer interface {
	Frame() int64
}

// Frame type represents a specific frame index within some media content
type Frame int64

// Frame gets the frame index for this frame. In this case, it just returns itself
func (f Frame) Frame() int64 {
	return int64(f)
}

// Equals checks if this framer is equal to the other
func (f Frame) Equals(other Framer) bool {
	return other.Frame() == f.Frame()
}

// Add adds another framer instance to this frame
func (f Frame) Add(other Framer) Frame {
	return Frame(int64(f) + other.Frame())
}
