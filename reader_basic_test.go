package gocsv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type model struct {
	Field1 string `csv:"header1"`
	Field2 string `csv:"header2"`
}

func TestReaderNilStream(t *testing.T) {
	assert.Panics(t, func() {
		r, err := NewReader[model](nil)
		assert.Nil(t, err)
		_, _ = r.Next()
	})
}
