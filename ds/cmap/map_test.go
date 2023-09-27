package cmap

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	cmap := New[string](32)

	fmt.Println(cmap)
}
