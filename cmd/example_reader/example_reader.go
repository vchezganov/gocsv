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
