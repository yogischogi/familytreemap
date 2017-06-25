// Package familytreemap calculates relative population frequencies
// from Family Tree DNA projects.
package main

import (
	"flag"
	"fmt"
	"os"
	"sort"

	"github.com/yogischogi/familytreemap/ftdna"
)

func main() {
	var (
		in       = flag.String("in", "", "Filename of input table.")
		out      = flag.String("out", "", "Filename for results in CSV format.")
		col      = flag.Int("col", 3, "Column number that contains country information.")
		totalsin = flag.String("totalsin", "", "Totals in: Number of testers from each country.")
		sumuk    = flag.Bool("sumuk", false, "Sum UK: Adds the number of testers from England, Wales, Scotland and Northern Ireland to United Kingdom.")
		statout  = flag.String("statout", "", "Filename for elaborate statistical information.")
	)
	flag.Parse()

	var finalFrequencies ftdna.Frequencies

	// Read countries from Family Tree DNA project spreadsheet.
	countries, err := ftdna.ReadCountriesFromCSV(*in, *col-1)
	if err != nil {
		fmt.Printf("Error reading project table, %v.\n", err)
		os.Exit(1)
	}

	countryFreqs := countries.Frequencies()

	// Add testers from England, Wales, Scotland and N. Ireland to United Kingdom.
	if *sumuk == true {
		countryFreqs.SumUKTesters()
	}

	// Calculate relative frequencies.
	if *totalsin != "" {
		// Read the total number of testers from each country.
		var totalTesters map[string]float32
		var err error
		if *totalsin != "" {
			totalTesters, err = ftdna.ReadCountryTestersFromCSV(*totalsin)
			if err != nil {
				fmt.Printf("Error reading the number of testers from file, %v.\n", err)
				os.Exit(1)
			}
		}

		// Add testers from England, Wales, Scotland and N. Ireland to United Kingdom.
		if *sumuk == true {
			totalTesters["United Kingdom"] += totalTesters["England"] +
				totalTesters["Wales"] +
				totalTesters["Scotland"] +
				totalTesters["Northern Ireland"]
		}

		// Write elaborate statistical information.
		if *statout != "" {
			err = ftdna.WriteStatisticsAsCSV(*statout, countryFreqs, totalTesters)
			if err != nil {
				fmt.Printf("Error writing statistics to file, %v.\n", err)
			}
		}

		// Calculate relative frequencies in percent.
		relFreqs := make([]ftdna.Frequency, 0)
		for _, freq := range countryFreqs {
			total := totalTesters[freq.Country]
			if total > 0 {
				rel := freq.Persons / total * 100
				relFreqs = append(relFreqs, ftdna.Frequency{Country: freq.Country, Persons: rel})
			}
		}
		finalFrequencies = relFreqs
	} else {
		finalFrequencies = countryFreqs
	}

	if *out != "" {
		sort.Stable(sort.Reverse(&finalFrequencies))
		// Write results to file.
		err = finalFrequencies.WriteCSV(*out)
		if err != nil {
			fmt.Printf("Error writing result to file in CSV format, %v.\r\n", err)
			os.Exit(1)
		}
	}
}
