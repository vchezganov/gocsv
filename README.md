### gocsv
Go package for parsing CSV records into structs. Currently, it supports only
the following types and its references:
- `int`, `int8`, `int16`, `int32`, `int64`
- `uint`, `uint8`, `uint16`, `uint32`, `uint64`
- `float32`, `float64`
- `string`
- 
In addition, you may provide own function to be used for parsing values. The function should accept `string` parameter and
return `error` if there are any errors when parsing.

## Examples

CSV records could be iterated using `Next()` method that return the model or an error if any:
```go
package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/vchezganov/gocsv"
)

type Person struct {
	Name string `csv:"name"`
	Age  uint   `csv:"age"`
	ID   string `csv:"pass,ParseID"`
}

func (p *Person) ParseID(value string) error {
	p.ID = fmt.Sprintf("ABC-%s", value)
	return nil
}

func main() {
	stringReader := strings.NewReader("age,pass,name\n32,12345,Vitaly\n45,54321,Alexey")

	reader, err := gocsv.NewReader[Person](stringReader)
	if err != nil {
		panic(err)
	}

	for {
		model, err := reader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			continue
		}

		fmt.Printf("Person: %v\n", model)
	}
}
```

Starting from Go 1.23 the cycle could be rewritten using simple `for` loop:
```go
for model, err := reader.Iter() {
	if err != nil {
			continue
	}
	
	fmt.Printf("Person: %v\n", model)
} 
```

`csv.Reader` is accessible through `CSVReader` field for specific CSV parsing options:
```go
reader, err := gocsv.NewReader[Person](stringReader)
if err != nil {
    panic(err)
}

reader.CSVReader.Comma = '#'
```

You may use `Marshaller` directly for complex cases:
```go
package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/vchezganov/gocsv"
)

type Person struct {
	Name     string `csv:"name"`
	Age      int    `csv:"age"`
	Location string `csv:"city"`
	ID       int    `csv:"passport,ParseID"`
}

func (p *Person) ParseID(value string) error {
	s, err := strconv.Atoi(value)
	if err != nil {
		return err
	}

	if 10000 <= s && s <= 99999 {
		p.ID = s
		return nil
	}

	return errors.New("ID is not valid")
}

func main() {
	s := strings.NewReader("name,age,city,passport\nVitaly,25,Bonn,10000")
	csvReader := csv.NewReader(s)
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

	fmt.Printf("Person: %v\n", *model)
}
```


### Next steps
- Support composition
```go
type A struct {
	ID        int    `csv:"id"`
	Timestamp string `csv:"timestamp"`
}

// id, timestamp, name
type B struct {
	*A
	Name int `csv:"name"`
}

// prefix_id, prefix_timestamp, name
type C struct {
	Base *A  `csv:"prefix"`
	Name int `csv:"name"`
}
```
- Parsing function to accept not only `string` but `int`, `float`, etc.
- Converting structs into CSV dictionary, slice or string