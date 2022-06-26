package gocsv

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type model struct {
	Field1 string `csv:"header1"`
	Field2 string `csv:"header2"`
}

func TestReaderNilStream(t *testing.T) {
	assert.Panics(t, func() {
		_, _ = NewReader[model](nil)
	})
}

func TestReaderEmptyFile(t *testing.T) {
	s := strings.NewReader("")
	_, err := NewReader[model](s)
	assert.NotNil(t, err)
}

func TestReaderNextEOF(t *testing.T) {
	s := strings.NewReader("header1,header2")
	r, err := NewReader[model](s)
	assert.Nil(t, err)

	_, err = r.Next()
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, io.EOF)
}

func TestReaderNext(t *testing.T) {
	s := strings.NewReader("header1,header2\nvalue1,value2")
	r, err := NewReader[model](s)
	assert.Nil(t, err)

	m, err := r.Next()
	assert.Nil(t, err)
	assert.Equal(t, m, &model{"value1", "value2"})

	_, err = r.Next()
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, io.EOF)
}

func TestReaderWithHeaders(t *testing.T) {
	s := strings.NewReader("value1,value2,value3")
	r, err := NewReaderWithHeaders[model](s, []string{"header2", "header3", "header1"})
	assert.Nil(t, err)

	m, err := r.Next()
	assert.Nil(t, err)
	assert.Equal(t, m, &model{"value3", "value1"})

	_, err = r.Next()
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, io.EOF)
}
