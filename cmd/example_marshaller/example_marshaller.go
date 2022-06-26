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
