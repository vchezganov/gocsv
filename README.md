### gocsv
Go package for parsing CSV records into structs. Currently, it supports only
the following types and its references:
- `int`, `int8`, `int16`, `int32`, `int64`
- `uint`, `uint8`, `uint16`, `uint32`, `uint64`
- `float32`, `float64`
- `string`

### Example
```go
package main

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/vchezganov/gocsv"
)

type Person struct {
	Name     string `csv:"name"`
	Age      int    `csv:"age"`
	Location string `csv:"city"`
}

func main() {
	f, err := os.Open("example.csv")
	if err != nil {
		panic(err)
	}

	csvReader := csv.NewReader(f)
	headers, err := csvReader.Read()
	if err != nil {
		panic(err)
	}

	marshaller, err := gocsv.NewMarshaller(headers)
	if err != nil {
		panic(err)
	}

	model := new(Person)
	records, err := csvReader.Read()
	if err != nil {
		panic(err)
	}

	err = marshaller.Unmarshal(records, model)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Model: %v", model)
}
```