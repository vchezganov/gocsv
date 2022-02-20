package gocsv

import (
	"bytes"
	"encoding/csv"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func stringPtr(s string) *string {
	return &s
}

func TestMarshaler(t *testing.T) {
	type testModel struct {
		Age      int     `csv:"age"`
		Name     string  `csv:"name"`
		Location *string `csv:"city"`
	}

	s := `name,age,city
Andrey,30,Moscow
Foma,45,
Kirill,22,Berlin`
	buf := bytes.NewBufferString(s)
	csvReader := csv.NewReader(buf)
	headers, err := csvReader.Read()
	assert.Nil(t, err)

	csvUnmarshal, err := NewMarshaller(headers)
	assert.Nil(t, err)

	expectedResults := []testModel{
		{Name: "Andrey", Age: 30, Location: stringPtr("Moscow")},
		{Name: "Foma", Age: 45, Location: nil},
		{Name: "Kirill", Age: 22, Location: stringPtr("Berlin")},
	}

	i := 0
	for ; ; i++ {
		records, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		assert.Nil(t, err)

		m := new(testModel)
		err = csvUnmarshal.Unmarshal(records, m)
		assert.Nil(t, err)
		assert.Equal(t, expectedResults[i], *m)
	}

	assert.Equal(t, i, len(expectedResults))
}

func TestMarshalerIntSizeError(t *testing.T) {
	type testModel struct {
		Number int8 `csv:"number"`
	}

	csvUnmarshal, err := NewMarshaller([]string{"number"})
	assert.Nil(t, err)

	m := new(testModel)
	err = csvUnmarshal.Unmarshal([]string{"512"}, m)
	assert.NotNil(t, err)
}
