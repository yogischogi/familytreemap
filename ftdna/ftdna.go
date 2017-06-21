// Package ftdna contains classes and methods for Family Tree DNA project data.
package ftdna

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// List of countries from a Family Tree DNA project spreadsheet.
type Countries []string

// ReadCountriesFromCSV reads a list of country names from a
// Family Tree DNA project spreadsheet.
// The spreadsheet must be in CSV format. countryCol is the number
// of the column that contains the countries, starting with 0.
func ReadCountriesFromCSV(filename string, countryCol int) (Countries, error) {
	infile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer infile.Close()

	// Read all CSV records from file.
	csvReader := csv.NewReader(infile)
	rows, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	// Add valid country entries to the result.
	result := make([]string, 0, len(rows))
	for _, row := range rows {
		if len(row) > countryCol {
			country := strings.TrimSpace(row[countryCol])
			if len(country) > 0 {
				result = append(result, country)
			}
		}
	}
	return result, nil
}

// Frequencies returns the number of testers from each country.
func (c *Countries) Frequencies() Frequencies {
	// Make a map that tells us how often a country appears in the list of countries.
	countryCount := make(map[string]float32)
	for _, name := range *c {
		if countryCount[name] > 0 {
			countryCount[name]++
		} else {
			countryCount[name] = 1
		}
	}

	// Convert map to list of frequencies.
	result := make([]Frequency, 0, len(countryCount))
	for country, persons := range countryCount {
		result = append(result, Frequency{Country: country, Persons: persons})
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
	_, err = writer.WriteString(fmt.Sprintf("%s,%s\r\n", "Country", "Testers"))
	if err != nil {
		return err
	}
	// Write rows.
	for _, freq := range *f {
		_, err = writer.WriteString(fmt.Sprintf("%s,%g\r\n", freq.Country, freq.Persons))
		if err != nil {
			return err
		}
	}
	err = writer.Flush()
	return err
}

// ReadCountryTestersFromCSV reads how many persons from which country have tested.
// The input is a CSV encoded file that contains the country names
// and the number of persons tested. Example:
//  Country,Testers
// 	Belarus,1000
//	Belgium,2000
//	Brazil,100
func ReadCountryTestersFromCSV(filename string) (map[string]float32, error) {
	infile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer infile.Close()

	// Read all CSV records from file.
	csvReader := csv.NewReader(infile)
	rows, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}
	// File must contain at least a header and one data row.
	if len(rows) < 2 {
		return nil, errors.New("not enough data in CSV file")
	}

	// Throw away header.
	rows = rows[1:]

	testers := make(map[string]float32)
	for _, row := range rows {
		country := row[0]
		cTesters, err := strconv.ParseFloat(row[1], 32)
		if err != nil {
			return nil, err
		}
		testers[country] = float32(cTesters)
	}
	return testers, nil
}
