// Package ftdna contains classes and methods for Family Tree DNA project data.
package ftdna

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Table contains the data of a Family Tree DNA project.
type Table struct {
	// countryCol is the table column number that contains the country names starting with 0.
	countryCol int
	// records is a table containing genetic results.
	records [][]string
}

// NewTableFromFile creates a new Table from a CSV encoded file.
// The file contains the same data as the project's genetic data
// as it appears on a web page.
func NewTableFromFile(filename string, countryCol int) (*Table, error) {
	infile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer infile.Close()

	// Read all CSV records from file.
	csvReader := csv.NewReader(infile)
	rawData, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	// Remove all rows that are too short to contain country information.
	records := make([][]string, 0, 0)
	for _, row := range rawData {
		if len(row) > countryCol {
			records = append(records, row)
		}
	}
	result := Table{countryCol: countryCol, records: records}
	return &result, nil
}

// FrequenciesOf calculates population frequencies for the given countries.
// countries contains a map of country names and the total number of persons
// who have been tested. For a country to be included in the evaluation
// the number of tested persons must be at least 1.
func (t *Table) FrequenciesOf(countries map[string]float32) Frequencies {
	// Create an empty map that contains only the countries we have data for.
	results := make(map[string]float32)
	for name, _ := range countries {
		results[name] = 0
	}

	// Calculate frequencies.
	for _, row := range t.records {
		// If the country name exists increase it's frequency value.
		countryName := strings.TrimSpace(row[t.countryCol])
		if countries[countryName] > 0 {
			results[countryName]++
		}
	}

	// Convert map to list of frequencies.
	result := make([]Frequency, 0, len(countries))
	for country, persons := range results {
		if persons > 0 {
			result = append(result, Frequency{Country: country, Persons: persons})
		}
	}
	return result
}

// Frequency shows how many persons from a country have been tested.
type Frequency struct {
	Country string
	Persons float32
}

// Frequencies is a list of Frequency that satisfies the sort.Interface.
type Frequencies []Frequency

func (f *Frequencies) Len() int {
	return len(*f)
}

// SumUKTesters adds the number of testers from the United Kingdom,
// England, Northern Ireland, Wales and Scotland together as United Kingdom.
func (f *Frequencies) SumUKTesters() {
	ukIdx := -1
	var sum float32 = 0
	for i, freq := range *f {
		switch {
		case freq.Country == "United Kingdom":
			sum += freq.Persons
			ukIdx = i
		case freq.Country == "England" ||
			freq.Country == "Wales" ||
			freq.Country == "Scotland" ||
			freq.Country == "Northern Ireland":
			sum += freq.Persons
		}
	}
	if ukIdx != 0 {
		(*f)[ukIdx].Persons = sum
	}
}

func (f *Frequencies) Less(i, j int) bool {
	if (*f)[i].Persons < (*f)[j].Persons {
		return true
	} else {
		return false
	}
}

func (f *Frequencies) Swap(i, j int) {
	(*f)[i], (*f)[j] = (*f)[j], (*f)[i]
}

// WriteCSV writes the frequencies to a file as comma separated values.
// It adds an header containing the captions "Location" and "Value".
func (f *Frequencies) WriteCSV(filename string) error {
	// Open file.
	outfile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer outfile.Close()

	// Write header.
	writer := bufio.NewWriter(outfile)
	_, err = writer.WriteString(fmt.Sprintf("%s,%s\r\n", "Location", "Value"))
	if err != nil {
		return err
	}
	// Write rows.
	for _, freq := range *f {
		_, err = writer.WriteString(fmt.Sprintf("%s,%f\r\n", freq.Country, freq.Persons))
		if err != nil {
			return err
		}
	}
	err = writer.Flush()
	return err
}

// ReadCountryTesters reads how many persons from which country have tested.
// The input is a JSON encoded file that contains the country names
// and the number of persons tested. Example:
// 	{"Belarus":1000,
//	"Belgium":2000,
//	"Brazil":100}
func ReadCountryTesters(filename string) (map[string]float32, error) {
	testers := make(map[string]float32)
	infile, err := os.Open(filename)
	if err != nil {
		return testers, err
	}
	defer infile.Close()

	decoder := json.NewDecoder(infile)
	err = decoder.Decode(&testers)
	if err != nil {
		return testers, err
	}
	return testers, nil
}
