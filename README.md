### gocsv
Go package for parsing CSV records into structs. Currently, it supports only
the following types and its references:
- `int`, `int8`, `int16`, `int32`, `int64`
- `uint`, `uint8`, `uint16`, `uint32`, `uint64`
- `float32`, `float64`
- `string`

### Example
```go
type Person struct {
	Name     string `csv:"name"`
	Age      int    `csv:"age"`
	Location string `csv:"city"`
}

...

f, _ := os.Open("...")
csvReader := csv.NewReader(f)
headers, _ := csvReader.Read()
marshaller, _ := NewMarshaller(headers)

model := new(Person)
records, _ := csvReader.Read()
_ = marshaller.Unmarshal(records, model)
```