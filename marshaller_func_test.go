package gocsv

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testParseFuncModel struct {
	Operation      string `csv:"operation,ParseOperation"`
	ErrorReturn    string `csv:"error_return,ParseErrorReturn"`
	Unexported     string `csv:"unexported,parseUnexported"`
	NoReturn       string `csv:"no_return,ParseNoReturn"`
	NotErrorReturn string `csv:"not_error_return,ParseNotErrorReturn"`
	NotStringArg   int    `csv:"not_string_arg,ParseNotStringArg"`
	NotOneArg      string `csv:"not_one_arg,ParseNotOneArg"`
}

func (m *testParseFuncModel) ParseOperation(value string) error {
	switch value {
	case "+":
		m.Operation = "add"

	case "-":
		m.Operation = "minus"

	case "*":
		m.Operation = "multiply"

	case "/":
		m.Operation = "divide"

	default:
		m.Operation = "add"
	}

	return nil
}

func (m *testParseFuncModel) ParseErrorReturn(_ string) error {
	return errors.New("error")
}

func (m *testParseFuncModel) parseUnexported(_ string) error {
	return nil
}

func (m *testParseFuncModel) ParseNoReturn(_ string) {
}

func (m *testParseFuncModel) ParseNotErrorReturn(_ string) string {
	return "error"
}

func (m *testParseFuncModel) ParseNotStringArg(_ int) error {
	return nil
}

func (m *testParseFuncModel) ParseNotOneArg(_, _ string) error {
	return nil
}

func TestMarshalerParserFunc(t *testing.T) {
	csvUnmarshal, err := NewMarshaller([]string{"operation"})
	assert.Nil(t, err)

	m := new(testParseFuncModel)
	err = csvUnmarshal.Unmarshal([]string{"+"}, m)
	assert.Nil(t, err)
	assert.Equal(t, "add", m.Operation)
}

func TestMarshalerParserErrorReturnFunc(t *testing.T) {
	csvUnmarshal, err := NewMarshaller([]string{"error_return"})
	assert.Nil(t, err)

	m := new(testParseFuncModel)
	err = csvUnmarshal.Unmarshal([]string{"+"}, m)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "error")
}

func TestMarshalerParserUnexportedFunc(t *testing.T) {
	csvUnmarshal, err := NewMarshaller([]string{"unexported"})
	assert.Nil(t, err)

	m := new(testParseFuncModel)
	err = csvUnmarshal.Unmarshal([]string{"+"}, m)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "cannot find")
}

func TestMarshalerParserNoReturnFunc(t *testing.T) {
	csvUnmarshal, err := NewMarshaller([]string{"no_return"})
	assert.Nil(t, err)

	m := new(testParseFuncModel)
	err = csvUnmarshal.Unmarshal([]string{"+"}, m)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "should return only one value")
}

func TestMarshalerParserNotErrorReturnFunc(t *testing.T) {
	csvUnmarshal, err := NewMarshaller([]string{"not_error_return"})
	assert.Nil(t, err)

	m := new(testParseFuncModel)
	err = csvUnmarshal.Unmarshal([]string{"+"}, m)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "should return only error type")
}

func TestMarshalerParserNotStringArg(t *testing.T) {
	csvUnmarshal, err := NewMarshaller([]string{"not_string_arg"})
	assert.Nil(t, err)

	m := new(testParseFuncModel)
	err = csvUnmarshal.Unmarshal([]string{"+"}, m)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "should accept only string type argument")
}

func TestMarshalerParserNotOneArg(t *testing.T) {
	csvUnmarshal, err := NewMarshaller([]string{"not_one_arg"})
	assert.Nil(t, err)

	m := new(testParseFuncModel)
	err = csvUnmarshal.Unmarshal([]string{"+"}, m)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "should accept only one argument")
}
