package gocsv

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReaderNextEmptyFile(t *testing.T) {
	s := strings.NewReader("")
	r, err := NewReader[model](s)
	assert.Nil(t, err)
	_, err = r.Next()
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, io.EOF)
}

func TestReaderNextOnlyHeader(t *testing.T) {
	s := strings.NewReader("header1,header2")
	r, err := NewReader[model](s)
	assert.Nil(t, err)

	_, err = r.Next()
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, io.EOF)
}

func TestReaderNext(t *testing.T) {
	s := strings.NewReader("header1,header2\nvalue1,value2\nvalue3,value4")
	r, err := NewReader[model](s)
	assert.Nil(t, err)

	m, err := r.Next()
	assert.Nil(t, err)
	assert.Equal(t, m, model{"value1", "value2"})

	m, err = r.Next()
	assert.Nil(t, err)
	assert.Equal(t, m, model{"value3", "value4"})

	_, err = r.Next()
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, io.EOF)
}

func TestReaderNextWithHeaders(t *testing.T) {
	s := strings.NewReader("value1,value2,value3")
	r, err := NewReaderWithHeaders[model](s, []string{"header2", "header3", "header1"})
	assert.Nil(t, err)

	m, err := r.Next()
	assert.Nil(t, err)
	assert.Equal(t, m, model{"value3", "value1"})

	_, err = r.Next()
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, io.EOF)
}

func TestReaderNextReferenceError(t *testing.T) {
	s := strings.NewReader("header1,header2\nvalue1,value2")
	r, err := NewReader[*model](s)
	assert.Nil(t, err)

	_, err = r.Next()
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, ErrNotStruct)
}

func TestReaderNextHashSeparator(t *testing.T) {
	s := strings.NewReader("header1#header2\nvalue1#value2\nvalue3#value4")
	r, err := NewReader[model](s)
	assert.Nil(t, err)

	r.CSVReader.Comma = '#'

	m, err := r.Next()
	assert.Nil(t, err)
	assert.Equal(t, m, model{"value1", "value2"})

	m, err = r.Next()
	assert.Nil(t, err)
	assert.Equal(t, m, model{"value3", "value4"})

	_, err = r.Next()
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, io.EOF)
}
