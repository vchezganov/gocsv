package gocsv

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

type marshaller struct {
	mapHeaders map[string]uint
}

func (p *marshaller) Unmarshal(values []string, model interface{}) error {
	reflectValue := reflect.ValueOf(model)
	if reflectValue.Kind() != reflect.Ptr || reflectValue.IsNil() {
		return errors.New("model must be a non-nil pointer")
	}

	reflectValue = reflectValue.Elem()
	reflectType := reflectValue.Type()
	if reflectType.Kind() != reflect.Struct {
		return errors.New("model must be a struct")
	}

	numFields := reflectType.NumField()
	for i := 0; i < numFields; i++ {
		reflectFieldType := reflectType.Field(i)
		header, ok := reflectFieldType.Tag.Lookup("csv")
		if !ok {
			continue
		}

		index, ok := p.mapHeaders[header]
		if !ok {
			continue
		}

		value := values[index]

		reflectFieldValue := reflectValue.Field(i)
		if reflectFieldValue.Kind() == reflect.Ptr {
			if value == "" {
				continue
			}

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
				return fmt.Errorf("cannot parse \"%s\" in \"%s\" column: %w", value, header, err)
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
				return fmt.Errorf("cannot parse \"%s\" in \"%s\" column: %w", value, header, err)
			}
			reflectFieldValue.SetUint(v)

		case
			reflect.Float32,
			reflect.Float64:
			v, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return fmt.Errorf("cannot parse \"%s\" in \"%s\" column: %w", value, header, err)
			}
			reflectFieldValue.SetFloat(v)

		case reflect.String:
			reflectFieldValue.SetString(value)

		default:
			return fmt.Errorf("unsupported field type: %s", k)
		}
	}

	return nil
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
