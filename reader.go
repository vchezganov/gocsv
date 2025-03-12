package gocsv

import (
	"encoding/csv"
	"io"
	"iter"
)

// Reader is an object to read and parse CSV files into Go objects.
type Reader[T any] struct {
	CSVReader  *csv.Reader
	headers    []string
	marshaller Marshaller
}

func (r *Reader[T]) init() error {
	if len(r.headers) == 0 {
		firstValues, err := r.CSVReader.Read()
		if err != nil {
			return err
		}

		r.headers = firstValues
	}

	csvMarshaller, err := NewMarshaller(r.headers)
	if err != nil {
		return err
	}

	r.marshaller = csvMarshaller
	return nil
}

func (r *Reader[T]) Iter() iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		for i := 0; true; i++ {
			model, err := r.Next()
			if err == io.EOF {
				break
			}

			if !yield(model, err) {
				break
			}
		}
	}
}

// Next method returns next instance or error if any.
// io.EOF error means that there are no more records.
func (r *Reader[T]) Next() (T, error) {
	var model T
	err := r.init()
	if err != nil {
		return model, err
	}

	records, err := r.CSVReader.Read()
	if err != nil {
		return model, err
	}

	//model := new(T)
	err = r.marshaller.Unmarshal(records, &model)
	if err != nil {
		return model, err
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
	return &Reader[T]{
		CSVReader: csvReader,
		headers:   headers,
	}, nil
}
