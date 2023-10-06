package typedef

// A DataBytes holds an immutable view of bytes.
type DataBytes struct {
	B []byte
}

// Len returns the view's length
func (v DataBytes) Len() int {
	return len(v.B)
}

// ByteSlice returns a copy of the data as a byte slice.
func (v DataBytes) ByteSlice() []byte {
	return CloneBytes(v.B)
}

// String returns the data as a string, making a copy if necessary.
func (v DataBytes) String() string {
	return string(v.B)
}

func CloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
