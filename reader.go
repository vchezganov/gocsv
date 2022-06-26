package gocsv

import (
	"encoding/csv"
	"io"
)

// Reader is an object to read and parse CSV files into Go objects.
type Reader[T any] struct {
	CSVReader  *csv.Reader
	marshaller Marshaller
}

// Next method returns next instance or error if any.
// io.EOF error means that there are no more records.
func (r *Reader[T]) Next() (*T, error) {
	records, err := r.CSVReader.Read()
	if err != nil {
		return nil, err
	}

	model := new(T)
	err = r.marshaller.Unmarshal(records, model)
	if err != nil {
		return nil, err
	}

	return model, nil
}

// NewReader creates an instance of gocsv.Reader from io.Reader.
// The first row values are considered as headers.
// If it is required to set headers manually, you might use NewReaderWithHeaders method.
func NewReader[T any](csvFile io.Reader) (*Reader[T], error) {
	return NewReaderWithHeaders[T](csvFile, nil)
}

// NewReaderWithHeaders creates an instance of gocsv.Reader from io.Reader and
// defined headers. If no headers are provider, the first row values are considered as headers.
func NewReaderWithHeaders[T any](csvFile io.Reader, headers []string) (*Reader[T], error) {
	csvReader := csv.NewReader(csvFile)
	if len(headers) == 0 {
		firstValues, err := csvReader.Read()
		if err != nil {
			return nil, err
		}

		headers = firstValues
	}

	csvMarshaller, err := NewMarshaller(headers)
	if err != nil {
		return nil, err
	}

	return &Reader[T]{
		CSVReader:  csvReader,
		marshaller: csvMarshaller,
	}, nil
}
