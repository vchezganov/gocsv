package gocsv

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type marshaller struct {
	mapHeaders map[string]uint
}

func (p *marshaller) Unmarshal(values []string, model interface{}) (err error) {
	reflectRef := reflect.ValueOf(model)
	if reflectRef.Kind() != reflect.Ptr || reflectRef.IsNil() {
		return errors.New("model must be a non-nil pointer")
	}

	reflectElem := reflectRef.Elem()
	reflectType := reflectElem.Type()
	if reflectType.Kind() != reflect.Struct {
		return errors.New("model must be a struct")
	}

	numFields := reflectType.NumField()
	for fieldIndex := 0; fieldIndex < numFields; fieldIndex++ {
		// Parsing field tag
		reflectFieldType := reflectType.Field(fieldIndex)
		csvTag, ok := reflectFieldType.Tag.Lookup("csv")
		if !ok {
			continue
		}

		splitTag := strings.SplitN(csvTag, ",", 2)
		header := splitTag[0]
		parser := ""
		if len(splitTag) == 2 {
			parser = splitTag[1]
		}

		index, ok := p.mapHeaders[header]
		if !ok {
			continue
		}

		// Parsing value
		value := values[index]

		if parser == "" {
			err = p.parseValue(reflectElem, fieldIndex, value)
		} else {
			err = p.parseFunc(reflectRef, parser, value)
		}

		if err != nil {
			err = fmt.Errorf("cannot parse \"%s\" in \"%s\" column: %w", value, header, err)
			return
		}
	}

	return
}

func (p *marshaller) parseValue(reflectElem reflect.Value, fieldIndex int, value string) error {
	reflectFieldValue := reflectElem.Field(fieldIndex)
	if reflectFieldValue.Kind() == reflect.Ptr {
		// Doing nothing if value is empty and field is a pointer
		if value == "" {
			return nil
		}

		// Creating an instance if field is null
		if reflectFieldValue.IsNil() {
			reflectPtrValue := reflect.New(reflectFieldValue.Type().Elem())
			reflectFieldValue.Set(reflectPtrValue)
		}

		reflectFieldValue = reflectFieldValue.Elem()
	}

	switch k := reflectFieldValue.Kind(); k {
	case
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		size := int(reflectFieldValue.Type().Size())
		v, err := strconv.ParseInt(value, 10, 8*size)
		if err != nil {
			return err
		}
		reflectFieldValue.SetInt(v)

	case
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		size := int(reflectFieldValue.Type().Size())
		v, err := strconv.ParseUint(value, 10, 8*size)
		if err != nil {
			return err
		}
		reflectFieldValue.SetUint(v)

	case
		reflect.Float32,
		reflect.Float64:
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		reflectFieldValue.SetFloat(v)

	case reflect.String:
		reflectFieldValue.SetString(value)

	default:
		return fmt.Errorf("unsupported field type: %s", k)
	}

	return nil
}

func (p *marshaller) parseFunc(reflectRef reflect.Value, methodName, value string) error {
	reflectMethod := reflectRef.MethodByName(methodName)
	if !reflectMethod.IsValid() {
		return fmt.Errorf("cannot find \"%s\" method", methodName)
	}

	reflectMethodType := reflectMethod.Type()

	if reflectMethodType.NumIn() != 1 {
		return fmt.Errorf("method \"%s\" should accept only one argument", methodName)
	}

	reflectIn := reflectMethodType.In(0)
	if reflectIn.Kind() != reflect.String {
		return fmt.Errorf("method \"%s\" should accept only string type argument", methodName)
	}

	if reflectMethodType.NumOut() != 1 {
		return fmt.Errorf("method \"%s\" should return only one value", methodName)
	}

	reflectOut := reflectMethodType.Out(0)
	if !reflectOut.Implements(reflect.TypeOf((*error)(nil)).Elem()) {
		return fmt.Errorf("method \"%s\" should return only error type interface", methodName)
	}

	callResult := reflectMethod.Call([]reflect.Value{
		reflect.ValueOf(value),
	})

	reflectResult := callResult[0]
	if reflectResult.Interface() == nil {
		return nil
	}

	err, ok := callResult[0].Interface().(error)
	if !ok {
		err = fmt.Errorf("method \"%s\" should return only error type interface", methodName)
	}

	return err
}

func NewMarshaller(headers []string) (*marshaller, error) {
	if len(headers) == 0 {
		return nil, fmt.Errorf("no headers")
	}

	mapHeaders := make(map[string]uint, len(headers))
	for i, header := range headers {
		mapHeaders[header] = uint(i)
	}

	return &marshaller{
		mapHeaders: mapHeaders,
	}, nil
}
