package base

type String string

func (s String) Len() int {
	return len(s)
}

type Integer int64

func (s Integer) Len() int {
	return 8
}

type Float float64

func (f Float) Len() int {
	return 8
}
