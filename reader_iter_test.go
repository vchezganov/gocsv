package gocsv

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReaderIterEmptyFile(t *testing.T) {
	s := strings.NewReader("")
	r, err := NewReader[model](s)
	assert.Nil(t, err)

	i := 0
	for _, err = range r.Iter() {
		i++
	}

	assert.Nil(t, err)
	assert.Equal(t, 0, i)
}

func TestReaderIterOnlyHeader(t *testing.T) {
	s := strings.NewReader("header1,header2")
	r, err := NewReader[model](s)
	assert.Nil(t, err)

	i := 0
	for _, err = range r.Iter() {
		i++
	}

	assert.Nil(t, err)
	assert.Equal(t, 0, i)
}

func TestReaderIter(t *testing.T) {
	s := strings.NewReader("header1,header2\nvalue1,value2\nvalue3,value4")
	r, err := NewReader[model](s)
	assert.Nil(t, err)

	expectedValues := []model{
		{"value1", "value2"},
		{"value3", "value4"},
	}

	i := 0
	for m, err := range r.Iter() {
		assert.Nil(t, err)
		assert.Equal(t, m, expectedValues[i])
		i++
	}

	assert.Equal(t, len(expectedValues), i)
}

func TestReaderIterWithHeaders(t *testing.T) {
	s := strings.NewReader("value1,value2,value3\nvalue4,value5,value6")
	r, err := NewReaderWithHeaders[model](s, []string{"header2", "header3", "header1"})
	assert.Nil(t, err)

	expectedValues := []model{
		{"value3", "value1"},
		{"value6", "value4"},
	}

	i := 0
	for m, err := range r.Iter() {
		assert.Nil(t, err)
		assert.Equal(t, m, expectedValues[i])
		i++
	}

	assert.Equal(t, len(expectedValues), i)
}

func TestReaderIterError(t *testing.T) {
	s := strings.NewReader("header1,header2\nvalue1,value2,value3\nvalue4,value5")
	r, err := NewReader[model](s)
	assert.Nil(t, err)

	expectedValue := model{"value4", "value5"}

	for m, err := range r.Iter() {
		if err == nil {
			assert.Equal(t, m, expectedValue)
		}
	}
}

func TestReaderIterReferenceError(t *testing.T) {
	s := strings.NewReader("header1,header2\nvalue1,value2\nvalue3,value4")
	r, err := NewReader[*model](s)
	assert.Nil(t, err)

	i := 0
	for _, err = range r.Iter() {
		assert.NotNil(t, err)
		assert.ErrorIs(t, err, ErrNotStruct)
		i++
	}

	assert.Equal(t, 2, i)
}

func TestReaderIterHashSeparator(t *testing.T) {
	s := strings.NewReader("header1#header2\nvalue1#value2\nvalue3#value4")
	r, err := NewReader[model](s)
	assert.Nil(t, err)

	r.CSVReader.Comma = '#'

	expectedValues := []model{
		{"value1", "value2"},
		{"value3", "value4"},
	}

	i := 0
	for m, err := range r.Iter() {
		assert.Nil(t, err)
		assert.Equal(t, m, expectedValues[i])
		i++
	}

	assert.Equal(t, len(expectedValues), i)
}
